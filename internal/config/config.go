package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySqlUrn string
	OpenOB   OpenOB `json:"openOB"`
}

type OpenOB struct {
	Path     string `json:"path,optional"`
	UserName string `json:"userName,optional"`
	Password string `json:"password,optional"`
}
