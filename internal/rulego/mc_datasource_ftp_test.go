package rulego

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jlaffaye/ftp"

	"workflow/internal/datasource"
)

func TestFileServerNode_ProcessFile(t *testing.T) {
	// 创建测试文件
	testContent := []byte("test content")
	err := os.WriteFile("./testdata/test.txt", testContent, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	defer os.Remove("./testdata/test.txt")

	tests := []struct {
		name    string
		config  datasource.FileServerConfig
		action  string
		srcPath string
		dstPath string
		wantErr bool
	}{
		{
			name: "SFTP上传测试",
			config: datasource.FileServerConfig{
				Protocol: "sftp",
				Host:     "10.99.169.7",
				Port:     2233,
				Username: "beuser",
				Password: "Bepassword@123",
			},
			action:  "upload",
			srcPath: "./testdata/test.txt",
			dstPath: "/tmp/test.txt",
		},
		{
			name: "SFTP下载测试",
			config: datasource.FileServerConfig{
				Protocol: "sftp",
				Host:     "10.99.169.7",
				Port:     2233,
				Username: "beuser",
				Password: "Bepassword@123",
			},
			action:  "download",
			srcPath: "/tmp/test.txt",
			dstPath: "./testdata/test_download.txt",
		},
		{
			name: "SFTP删除测试",
			config: datasource.FileServerConfig{
				Protocol: "sftp",
				Host:     "10.99.169.7",
				Port:     2233,
				Username: "beuser",
				Password: "Bepassword@123",
			},
			action:  "delete",
			srcPath: "/tmp/test.txt",
		},
		{
			name: "FTP上传测试",
			config: datasource.FileServerConfig{
				Protocol: "ftp",
				Host:     "10.99.113.114",
				Port:     21,
				Username: "test",
				Password: "test",
			},
			action:  "upload",
			srcPath: "./testdata/test.txt",
			dstPath: "/test.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &FileServerNode{
				Config: FileServerNodeConfiguration{
					Type:     tt.config.Protocol,
					Action:   tt.action,
					DestPath: tt.dstPath,
				},
			}

			handler, err := node.createHandler(tt.config)
			if err != nil {
				t.Fatalf("创建处理器失败: %v", err)
			}
			defer handler.Close()

			var err2 error
			switch tt.action {
			case "upload":
				err2 = handler.Upload(tt.srcPath, tt.dstPath)
			case "download":
				err2 = handler.Download(tt.srcPath, tt.dstPath)
			case "delete":
				err2 = handler.Delete(tt.srcPath)
			}

			if (err2 != nil) != tt.wantErr {
				t.Errorf("执行%s操作失败: %v", tt.action, err2)
			}
		})
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
