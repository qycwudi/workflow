package datasource

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	goora "github.com/sijms/go-ora/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/ssh"

	"workflow/internal/enum"
)

type DataSourceConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

const (
	FtpProtocol  = "ftp"
	SftpProtocol = "sftp"
)

// FtpNodeConfiguration FTP节点配置
type FileServerConfig struct {
	Protocol string `json:"protocol"` // ftp 或 sftp
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckDataSourceClient(t enum.DBType, config string) error {
	// 根据类型区分检查方法
	switch t {
	case enum.MysqlType, enum.OracleType, enum.SqlServerType:
		return checkDatabaseConnection(t, config)
	case enum.FileServerType:
		return checkFileServerConnection(config)
	default:
		return errors.New("unknown data source type")
	}
}

// 检查数据库连接
func checkDatabaseConnection(t enum.DBType, config string) error {
	dsn := GenDataSourceDSN(t, config)
	var db *sql.DB
	var err error

	switch t {
	case enum.MysqlType:
		db, err = sql.Open("mysql", dsn)
	case enum.OracleType:
		db, err = sql.Open("oracle", dsn)
	case enum.SqlServerType:
		db, err = sql.Open("sqlserver", dsn)
	}

	if err != nil {
		return fmt.Errorf("connect to database failed: %v", err)
	}

	defer func() {
		_ = db.Close()
	}()

	if err = db.Ping(); err != nil {
		return fmt.Errorf("ping database failed: %v", err)
	}
	return nil
}

// 检查文件服务器连接
func checkFileServerConnection(config string) error {
	var serverConfig FileServerConfig
	if err := json.Unmarshal([]byte(config), &serverConfig); err != nil {
		return fmt.Errorf("parse file server config failed: %v", err)
	}

	switch serverConfig.Protocol {
	case FtpProtocol:
		return checkFtpConnection(serverConfig)
	case SftpProtocol:
		return checkSftpConnection(serverConfig)
	default:
		return fmt.Errorf("unsupported file server protocol: %s", serverConfig.Protocol)
	}
}

// 检查FTP连接
func checkFtpConnection(config FileServerConfig) error {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ftp.Dial(addr)
	if err != nil {
		return fmt.Errorf("connect to ftp server failed: %v", err)
	}
	defer client.Quit()

	if err = client.Login(config.Username, config.Password); err != nil {
		return fmt.Errorf("ftp login failed: %v", err)
	}
	return nil
}

// 检查SFTP连接
func checkSftpConnection(config FileServerConfig) error {
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("connect to sftp server failed: %v", err)
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("create sftp client failed: %v", err)
	}
	defer sftpClient.Close()

	return nil
}

func GenDataSourceDSN(t enum.DBType, config string) string {
	c := DataSourceConfig{}
	err := json.Unmarshal([]byte(config), &c)
	if err != nil {
		logx.Errorf("unmarshal datasource config failed, err:%v", err)
		return ""
	}

	var dsn string
	switch t {
	case enum.MysqlType:
		// {"host": "192.168.49.2", "port": 31426, "database": "wk", "password": "root", "user": "root"}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
	case enum.OracleType:
		dsn = goora.BuildUrl(c.Host, c.Port, c.Database, c.User, c.Password, nil)
	case enum.SqlServerType:
		dsn = fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s", c.Host, c.Port, c.User, c.Password, c.Database)
	}

	return dsn
}
