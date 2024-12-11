package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySqlUrn string
	Job      JobConfig `json:"Job"`
}

// 任务配置
type JobConfig struct {
	DatasourceClientCheck JobConfigDatasourceClientCheck `json:"DatasourceClientCheck"`
}

type JobConfigDatasourceClientCheck struct {
	JobBase
}

type JobBase struct {
	Enable bool   `json:"enable" default:"false"`
	Cron   string `json:"cron" default:"0"`
}
