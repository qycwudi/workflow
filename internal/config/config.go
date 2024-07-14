package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySqlDataSource           string
	RocketMQNameServerAddress string
}
