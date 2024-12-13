package datasource

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	goora "github.com/sijms/go-ora/v2"
	"github.com/zeromicro/go-zero/core/logx"

	"workflow/internal/enum"
)

type DataSourceConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func CheckDataSourceClient(t enum.DBType, config string) error {
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
	default:
		err = errors.New("unknown data source type")
	}

	if err != nil {
		return errors.New("connect to datasource failed")
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.Ping()
	if err != nil {
		return errors.New("connect to datasource failed")
	}
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
