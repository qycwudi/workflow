package asynq

import (
	"github.com/hibiken/asynq"
	"gogogo/internal/config"
	"log"
)

type AsynqTask struct {
	AsynqTaskClient *asynq.Client
}

var AsynqTaskContext AsynqTask

func NewAsynqClient(config config.RedisConfig) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: config.RedisAddr, Username: config.RedisUserName, Password: config.RedisPassword, DB: config.RedisDb})
}

func NewAsynqServer(config config.RedisConfig) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.RedisAddr, Username: config.RedisUserName, Password: config.RedisPassword, DB: config.RedisDb},
		asynq.Config{
			Concurrency: 5,
			Queues: map[string]int{
				// "critical": 3,
				"default": 5,
				// "low":      1,
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeLLMFeatureExtraction, HandleLlmTask)
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
