package rulego

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jlaffaye/ftp"
	"github.com/rulego/rulego/api/types"
)

func TestSftpNode_executeFtp(t *testing.T) {
	// 创建测试文件
	testContent := []byte("test content")
	err := os.WriteFile("./testdata/test.txt", testContent, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove("./testdata/test.txt")

	// 测试 SFTP
	msg := types.RuleMsg{
		Data: `{
			"action": "upload",
			"config": {
				"protocol": "sftp",
				"host": "10.99.169.7",
				"port": 2233,
				"username": "beuser",
				"password": "Bepassword@123"
			},
			"srcPath": "./testdata/test.txt",
			"destPath": "/tmp/test.txt"
		}`,
	}

	node := &FtpNode{}

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行SFTP操作失败: %v", err)
	}

	msg.Data = `{
		"action": "download",
		"config": {
			"protocol": "sftp",
			"host": "10.99.169.7",
			"port": 2233,
			"username": "beuser", 
			"password": "Bepassword@123"
		},
		"srcPath": "/tmp/test.txt",
		"destPath": "./testdata/test_download.txt"
	}`

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行SFTP下载操作失败: %v", err)
	}
	defer os.Remove("./testdata/test_download.txt")

	msg.Data = `{
		"action": "delete",
		"config": {
			"protocol": "sftp",
			"host": "10.99.169.7",
			"port": 2233,
			"username": "beuser",
			"password": "Bepassword@123"
		},
		"path": "/tmp/test.txt"
	}`

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行SFTP删除操作失败: %v", err)
	}
}

/*

docker pull stilliard/pure-ftpd

docker run -d \
    --name ftp-server \
    -e FTP_USER_NAME=ftpuser \
    -e FTP_USER_PASS=ftppassword \
    -e FTP_USER_HOME=/home/ftpuser \
    -v /Users/qiangyuecheng/images:/home/ftpuser \
    -p 21:21 -p 21000-21010:21000-21010 \
    stilliard/pure-ftpd

*/

func TestFtpNode_executeFtp(t *testing.T) {
	// 创建测试文件
	testContent := []byte("test content")
	err := os.WriteFile("./testdata/test.txt", testContent, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove("./testdata/test.txt")

	// 测试 FTP
	msg := types.RuleMsg{
		Data: `{"action":"upload","config":{"protocol":"ftp","host":"10.99.113.114","port":21,"username":"test","password":"test","passive":true},"srcPath":"./testdata/test.txt","destPath":"/test.txt"}`,
	}

	node := &FtpNode{}

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行FTP上传操作失败: %v", err)
	}

	msg.Data = `{
		"action": "download",
		"config": {
			"protocol": "ftp",
			"host": "10.99.113.114",
			"port": 21,
			"username": "test",
			"password": "test",
			"passive": true
		},
		"srcPath": "/test.txt",
		"destPath": "./testdata/test_download_ftp.txt"
	}`

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行FTP下载操作失败: %v", err)
	}
	defer os.Remove("./testdata/test_download_ftp.txt")

	msg.Data = `{
		"action": "delete",
		"config": {
			"protocol": "ftp",
			"host": "10.99.113.114",
			"port": 21,
			"username": "test",
			"password": "test",
			"passive": true
		},
		"path": "/test.txt"
	}`

	err = node.executeFtp(msg)
	if err != nil {
		t.Errorf("执行FTP删除操作失败: %v", err)
	}
}
func TestFtpNode_net_executeFtp(t *testing.T) {
	ftpServer := "10.99.113.114:21"
	username := "test"
	password := "test"

	// 连接到 FTP 服务器
	c, err := ftp.Dial(ftpServer)
	if err != nil {
		log.Fatalf("Failed to connect to FTP server: %v", err)
	}

	// 登录
	err = c.Login(username, password)
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	defer c.Logout()

	// 列出目录
	entries, err := c.List("")
	if err != nil {
		log.Fatalf("Failed to list directory: %v", err)
	}

	for _, entry := range entries {
		fmt.Println(entry.Name)
	}
}
