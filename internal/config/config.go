package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	ApiPort int
	rest.RestConf
	MySqlUrn string
	Job      []JobBase   `json:"job"`
	Redis    RedisConfig `json:"Redis"`
}

type RedisConfig struct {
	Host     string `json:"Host"`
	Password string `json:"Password"`
	DB       int    `json:"DB"`
}

type JobBase struct {
	Name   string `json:"name"`   // 任务名称
	Enable bool   `json:"enable"` // 是否启用
	Cron   string `json:"cron"`   // cron表达式
	Topic  string `json:"topic"`  // 任务主题
}
