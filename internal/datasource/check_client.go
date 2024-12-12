package datasource

import (
	"database/sql"
	"errors"
	"fmt"
	goora "github.com/sijms/go-ora/v2"
	"github.com/tidwall/gjson"
	"workflow/internal/enum"
)

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
	}

	if err != nil {
		return errors.New("connect to datasource failed")
	}
	defer func() { _ = db.Close() }()

	err = db.Ping()
	if err != nil {
		return errors.New("connect to datasource failed")
	}
	return nil
}

func GenDataSourceDSN(t enum.DBType, config string) string {
	host := gjson.Get(config, "host").String()
	port := gjson.Get(config, "port").Int()
	user := gjson.Get(config, "user").String()
	password := gjson.Get(config, "password").String()
	database := gjson.Get(config, "database").String()

	var dsn string
	switch t {
	case enum.MysqlType:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)
	case enum.OracleType:
		dsn = goora.BuildUrl(host, int(port), database, user, password, nil)
	case enum.SqlServerType:
		dsn = fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s", host, port, user, password, database)
	}

	return dsn
}
