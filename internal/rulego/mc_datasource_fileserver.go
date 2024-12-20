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

// FtpNode A plugin that flow sftp node
type FileServerNode struct {
	Config FileServerNodeConfiguration
}

// FtpAction FTP操作类型
type FileServerNodeConfiguration struct {
	Type         string `json:"type"`         // 数据源类型 ftp 或 sftp
	Action       string `json:"action"`       // upload/download/delete
	DatasourceId int64  `json:"datasourceId"` // 数据源 ID
	DestPath     string `json:"destPath"`     // 目标文件路径
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
func (n *FileServerNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	// 读取配置
	marshal, _ := json.Marshal(configuration)
	_ = json.Unmarshal(marshal, &n.Config)
	return nil
}

// OnMsg 处理消息
func (n *FileServerNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	source, err := RoleChain.svc.DatasourceModel.FindOne(context.Background(), n.Config.DatasourceId)
	if err != nil {
		ctx.TellFailure(msg, err)
		return
	}

	// 读取FTP配置
	var config datasource.FileServerConfig
	if err := json.Unmarshal([]byte(source.Config), &config); err != nil {
		ctx.TellFailure(msg, err)
		return
	}
	// 解析 data
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
		ctx.TellFailure(msg, err)
		return
	}
	// 生成临时文件路径
	var tmpPath string
	if tmp, ok := data["tmpPath"]; ok {
		tmpPath = tmp.(string)
	} else {
		tmpPath = filepath.Join(fmt.Sprintf("./file-server-%s", time.Now().Format("20060102150405")), filepath.Base(n.Config.DestPath))
		data["tmpPath"] = tmpPath
	}

	if err := n.executeFileServer(tmpPath, config, n.Config); err != nil {
		ctx.TellFailure(msg, err)
		return
	}
	ctx.TellSuccess(msg)
}

func (n *FileServerNode) Destroy() {
	// Do some cleanup work
}

// executeFileServer 执行文件服务器操作
func (n *FileServerNode) executeFileServer(tmpPath string, config datasource.FileServerConfig, configuration FileServerNodeConfiguration) error {
	switch config.Protocol {
	case datasource.SftpProtocol:
		return n.executeSftp(tmpPath, config, configuration)
	case datasource.FtpProtocol:
		return n.executeFtp(tmpPath, config, configuration)
	default:
		return fmt.Errorf("unsupported protocol: %s", config.Protocol)
	}
}

// executeFtpAction 执行FTP操作
func (n *FileServerNode) executeFtp(tmpPath string, config datasource.FileServerConfig, action FileServerNodeConfiguration) error {
	// 连接FTP服务器
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}
	defer client.Quit()

	// 登录
	if err := client.Login(config.Username, config.Password); err != nil {
		return err
	}

	// 根据操作类型执行不同的FTP操作
	switch action.Action {
	case "upload":
		return n.uploadFtpFile(client, tmpPath, action.DestPath)
	case "download":
		return n.downloadFtpFile(client, action.DestPath, tmpPath)
	case "delete":
		return client.Delete(action.DestPath)
	default:
		return fmt.Errorf("unsupported action: %s", action.Action)
	}
}

// executeSftp 执行SFTP操作
func (n *FileServerNode) executeSftp(tmpPath string, config datasource.FileServerConfig, configuration FileServerNodeConfiguration) error {
	// 配置SSH客户端
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// 连接SSH
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	// 创建SFTP客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// 根据操作类型执行不同的SFTP操作
	switch configuration.Action {
	case "upload":
		return n.uploadSftpFile(sftpClient, tmpPath, configuration.DestPath)
	case "download":
		return n.downloadSftpFile(sftpClient, tmpPath, configuration.DestPath)
	case "delete":
		return sftpClient.Remove(configuration.DestPath)
	default:
		return fmt.Errorf("unsupported action: %s", configuration.Action)
	}
}

// uploadFtpFile 上传文件(FTP)
func (n *FileServerNode) uploadFtpFile(client *ftp.ServerConn, srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	return client.Stor(destPath, srcFile)
}

// downloadFtpFile 下载文件(FTP)
func (n *FileServerNode) downloadFtpFile(client *ftp.ServerConn, srcPath, destPath string) error {
	resp, err := client.Retr(srcPath)
	if err != nil {
		return err
	}
	defer resp.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, resp)
	return err
}

// uploadSftpFile 上传文件(SFTP)
func (n *FileServerNode) uploadSftpFile(client *sftp.Client, srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := client.Create(destPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// downloadSftpFile 下载文件(SFTP)
func (n *FileServerNode) downloadSftpFile(client *sftp.Client, srcPath, destPath string) error {
	srcFile, err := client.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
