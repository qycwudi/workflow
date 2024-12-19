package rulego

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"github.com/rulego/rulego"
	"github.com/rulego/rulego/api/types"
	"github.com/rulego/rulego/utils/json"
	"golang.org/x/crypto/ssh"
)

// FtpNode A plugin that flow sftp node
type FtpNode struct{}

func init() {
	_ = rulego.Registry.Register(&FtpNode{})
}

func (n *FtpNode) Type() string {
	return Ftp
}
func (n *FtpNode) New() types.Node {
	return &FtpNode{}
}
func (n *FtpNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
	return nil
}

// OnMsg 处理消息
func (n *FtpNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) {
	if err := n.executeFtp(msg); err != nil {
		ctx.TellFailure(msg, err)
		return
	}
	ctx.TellSuccess(msg)
}

func (n *FtpNode) Destroy() {

	// Do some cleanup work
}

// FtpNodeConfiguration FTP节点配置
type FtpNodeConfiguration struct {
	Protocol string `json:"protocol"` // ftp 或 sftp
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// FtpAction FTP操作类型
type FtpAction struct {
	Action   string                 `json:"action"`   // upload/download/delete
	Config   map[string]interface{} `json:"config"`   // FTP配置
	SrcPath  string                 `json:"srcPath"`  // 源文件路径
	DestPath string                 `json:"destPath"` // 目标文件路径
	Path     string                 `json:"path"`     // 文件路径(用于删除)
}

// executeFtp 执行FTP操作
func (n *FtpNode) executeFtp(msg types.RuleMsg) error {
	// 解析消息数据
	var action FtpAction
	if err := json.Unmarshal([]byte(msg.Data), &action); err != nil {
		return err
	}

	// 解析FTP配置
	var config FtpNodeConfiguration
	configBytes, err := json.Marshal(action.Config)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return err
	}

	switch config.Protocol {
	case "sftp":
		return n.executeSftp(config, action)
	default:
		return n.executeFtpAction(config, action)
	}
}

// executeFtpAction 执行FTP操作
func (n *FtpNode) executeFtpAction(config FtpNodeConfiguration, action FtpAction) error {
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
		return n.uploadFtpFile(client, action.SrcPath, action.DestPath)
	case "download":
		return n.downloadFtpFile(client, action.SrcPath, action.DestPath)
	case "delete":
		return client.Delete(action.Path)
	default:
		return fmt.Errorf("unsupported action: %s", action.Action)
	}
}

// executeSftp 执行SFTP操作
func (n *FtpNode) executeSftp(config FtpNodeConfiguration, action FtpAction) error {
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
	switch action.Action {
	case "upload":
		return n.uploadSftpFile(sftpClient, action.SrcPath, action.DestPath)
	case "download":
		return n.downloadSftpFile(sftpClient, action.SrcPath, action.DestPath)
	case "delete":
		return sftpClient.Remove(action.Path)
	default:
		return fmt.Errorf("unsupported action: %s", action.Action)
	}
}

// uploadFtpFile 上传文件(FTP)
func (n *FtpNode) uploadFtpFile(client *ftp.ServerConn, srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	return client.Stor(destPath, srcFile)
}

// downloadFtpFile 下载文件(FTP)
func (n *FtpNode) downloadFtpFile(client *ftp.ServerConn, srcPath, destPath string) error {
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
func (n *FtpNode) uploadSftpFile(client *sftp.Client, srcPath, destPath string) error {
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
func (n *FtpNode) downloadSftpFile(client *sftp.Client, srcPath, destPath string) error {
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
