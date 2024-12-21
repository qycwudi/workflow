package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	ApiPort int
	rest.RestConf
	MySqlUrn string
	Job      JobConfig   `json:"Job"`
	Redis    RedisConfig `json:"Redis"`
}

type RedisConfig struct {
	Host     string `json:"Host"`
	Password string `json:"Password"`
	DB       int    `json:"DB"`
}

// 任务配置
type JobConfig struct {
	DatasourceClientCheck  JobConfigDatasourceClientCheck  `json:"DatasourceClientCheck"`
	DatasourceClientUpdate JobConfigDatasourceClientUpdate `json:"DatasourceClientUpdate"`
}

type JobConfigDatasourceClientCheck struct {
	JobBase
}

type JobConfigDatasourceClientUpdate struct {
	JobBase
}

type JobBase struct {
	Enable bool   `json:"enable" default:"false"`
	Cron   string `json:"cron" default:"0"`
}
