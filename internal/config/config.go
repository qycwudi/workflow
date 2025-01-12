package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	ApiPort int
	rest.RestConf
	MySqlUrn            string
	Job                 []JobBase   `json:"job"`
	Redis               RedisConfig `json:"Redis"`
	RuleServerLimitSize int         `json:"RuleServerLimitSize" default:"4"` // 规则链服务限制大小 单位 M
	RuleServerTrace     bool        `json:"RuleServerTrace" default:"true"`  // 是否开启规则链服务追踪
	Auth                struct {    // JWT 认证需要的密钥和过期时间配置
		AccessSecret string
		AccessExpire int64
	}
}

type RedisConfig struct {
	Host     string `json:"Host"`
	Password string `json:"Password"`
	DB       int    `json:"DB"`
}

type JobBase struct {
	Name   string `json:"name"`                  // 任务名称
	Enable bool   `json:"enable" default:"true"` // 是否启用
	Cron   string `json:"cron"`                  // cron表达式
}
