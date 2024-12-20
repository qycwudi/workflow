package rulego

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/json"
	"golang.org/x/crypto/ssh"

	"workflow/internal/datasource"
)

// FileServerNode 文件服务器节点
type FileServerNode struct {
	Config FileServerNodeConfiguration
}

// FileServerNodeConfiguration 节点配置
type FileServerNodeConfiguration struct {
	Type         string `json:"type"`         // 数据源类型 ftp 或 sftp
	Action       string `json:"action"`       // upload/download/delete
	DatasourceId int64  `json:"datasourceId"` // 数据源 ID
	DestPath     string `json:"destPath"`     // 目标文件路径
}

// FileServerHandler 文件服务器操作接口
type FileServerHandler interface {
	Upload(srcPath, destPath string) error
	Download(srcPath, destPath string) error
	Delete(path string) error
	Close() error
}

func init() {
	_ = rulego.Registry.Register(&FileServerNode{})
}

func (n *FileServerNode) Type() string {
	return FileServer
}

func (n *FileServerNode) New() types.Node {
	return &FileServerNode{}
}
func (n *FileServerNode) Destroy() {}

func (n *FileServerNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	marshal, _ := json.Marshal(configuration)
	return json.Unmarshal(marshal, &n.Config)
}

func (n *FileServerNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	source, err := RoleChain.svc.DatasourceModel.FindOne(context.Background(), n.Config.DatasourceId)
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	var config datasource.FileServerConfig
	if err := json.Unmarshal([]byte(source.Config), &config); err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	tmpPath := n.getTempPath(data)
	data["tmpPath"] = tmpPath

	if err := n.processFile(tmpPath, config); err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	marshal, _ := json.Marshal(data)
	msg.Data = string(marshal)
	ctx.TellSuccess(msg)
}

func (n *FileServerNode) getTempPath(data map[string]interface{}) string {
	if tmp, ok := data["tmpPath"]; ok {
		return tmp.(string)
	}
	return filepath.Join(fmt.Sprintf("./file-server-%s", time.Now().Format("20060102150405")), filepath.Base(n.Config.DestPath))
}

func (n *FileServerNode) processFile(tmpPath string, config datasource.FileServerConfig) error {
	handler, err := n.createHandler(config)
	if err != nil {
		return err
	}
	defer handler.Close()

	switch n.Config.Action {
	case "upload":
		return handler.Upload(tmpPath, n.Config.DestPath)
	case "download":
		return handler.Download(n.Config.DestPath, tmpPath)
	case "delete":
		return handler.Delete(n.Config.DestPath)
	default:
		return fmt.Errorf("unsupported action: %s", n.Config.Action)
	}
}

func (n *FileServerNode) createHandler(config datasource.FileServerConfig) (FileServerHandler, error) {
	switch config.Protocol {
	case datasource.FtpProtocol:
		return newFtpHandler(config)
	case datasource.SftpProtocol:
		return newSftpHandler(config)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", config.Protocol)
	}
}

// FtpHandler FTP处理器
type FtpHandler struct {
	client *ftp.ServerConn
}

func newFtpHandler(config datasource.FileServerConfig) (*FtpHandler, error) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	if err := client.Login(config.Username, config.Password); err != nil {
		client.Quit()
		return nil, err
	}

	return &FtpHandler{client: client}, nil
}

func (h *FtpHandler) Upload(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dirPath := filepath.Dir(destPath)
	if err := h.client.ChangeDir(dirPath); err != nil {
		if err := h.client.MakeDir(dirPath); err != nil {
			return err
		}
	}
	return h.client.Stor(destPath, srcFile)
}

func (h *FtpHandler) Download(srcPath, destPath string) error {
	resp, err := h.client.Retr(srcPath)
	if err != nil {
		return err
	}
	defer resp.Close()

	return saveFile(destPath, resp)
}

func (h *FtpHandler) Delete(path string) error {
	return h.client.Delete(path)
}

func (h *FtpHandler) Close() error {
	return h.client.Quit()
}

// SftpHandler SFTP处理器
type SftpHandler struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func newSftpHandler(config datasource.FileServerConfig) (*SftpHandler, error) {
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, err
	}

	return &SftpHandler{
		sshClient:  sshClient,
		sftpClient: sftpClient,
	}, nil
}

func (h *SftpHandler) Upload(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dirPath := filepath.Dir(destPath)
	if _, err := h.sftpClient.Stat(dirPath); err != nil {
		if err := h.sftpClient.MkdirAll(dirPath); err != nil {
			return err
		}
	}

	dstFile, err := h.sftpClient.Create(destPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (h *SftpHandler) Download(srcPath, destPath string) error {
	srcFile, err := h.sftpClient.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	return saveFile(destPath, srcFile)
}

func (h *SftpHandler) Delete(path string) error {
	return h.sftpClient.Remove(path)
}

func (h *SftpHandler) Close() error {
	h.sftpClient.Close()
	return h.sshClient.Close()
}

// 通用工具函数
func saveFile(destPath string, reader io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, reader)
	return err
}
