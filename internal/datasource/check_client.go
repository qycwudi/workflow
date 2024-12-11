package datasource

import (
	"database/sql"
	"errors"
)

func CheckDataSourceClient(t string, dsn string) error {
	if t == "mysql" {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return errors.New("connect to datasource failed")
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			return errors.New("connect to datasource failed")
		}
	}
	return nil
}
