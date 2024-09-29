package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySqlDataSource string
	OpenOB          OpenOB `json:"openOB"`
}

type OpenOB struct {
	Path     string `json:"path"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}
