package mq

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-clients/golang"
	"github.com/apache/rocketmq-clients/golang/credentials"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	// RockerMQConsumerAwaitDuration maximum waiting time for receive func
	RockerMQConsumerAwaitDuration = time.Second * 5

	// 并发度
	receiveConcurrency = 2

	// RockerMQConsumerMaxMessageNum maximum number of messages received at one time
	rockerConsumerMaxMessageNum int32 = 1

	// RockerMQConsumerInvisibleDuration invisibleDuration should > 20s
	rockerConsumerInvisibleDuration = time.Second * 20
)

type RocketMQConsumer struct {
	ocrConsumer golang.SimpleConsumer
	llmConsumer golang.SimpleConsumer
}

func NewRockerConsumer(ocrConsumer golang.SimpleConsumer, llmConsumer golang.SimpleConsumer) *RocketMQConsumer {
	return &RocketMQConsumer{ocrConsumer: ocrConsumer, llmConsumer: llmConsumer}
}

func (l *RocketMQConsumer) ReceiveOCRMessage() {
	logx.Info("start OCR consumer,listener topic:", OCR_V1_TOPIC)
	// Each Receive call will only select one broker queue to pop messages.
	// Enable multiple consumption goroutines to reduce message end-to-end latency.
	ch := make(chan struct{})
	wg := &sync.WaitGroup{}
	for i := 0; i < receiveConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ch:
					return
				default:
					mvs, _ := l.ocrConsumer.Receive(context.TODO(), rockerConsumerMaxMessageNum, rockerConsumerInvisibleDuration)
					// if err != nil {
					// 	fmt.Println("receive message error: " + err.Error())
					// }
					// ack message
					for _, mv := range mvs {
						logx.Info("consumption ocr message:", mv.GetMessageId())
						logx.Infof("body:%s \n", mv.GetBody())
						if err := l.ocrConsumer.Ack(context.TODO(), mv); err != nil {
							fmt.Println("ack message error: " + err.Error())
						}
					}
				}
			}
		}()
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	// wait for exit
	<-exit
	close(ch)
	wg.Wait()
	println("receiveOCRMessage close")
}

func (l *RocketMQConsumer) ReceiveLLMMessage() {
	logx.Info("start LLM consumer,listener topic:", LLM_TRAIT_EXTRACT_V1_TOPIC)
	// Each Receive call will only select one broker queue to pop messages.
	// Enable multiple consumption goroutines to reduce message end-to-end latency.
	ch := make(chan struct{})
	wg := &sync.WaitGroup{}
	for i := 0; i < receiveConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ch:
					return
				default:
					mvs, _ := l.llmConsumer.Receive(context.TODO(), rockerConsumerMaxMessageNum, rockerConsumerInvisibleDuration)
					// if err != nil {
					// 	fmt.Println("receive message error: " + err.Error())
					// }
					// ack message
					for _, mv := range mvs {
						logx.Info("consumption llm message:", mv.GetMessageId())
						time.Sleep(time.Second * 5)
						logx.Infof("message body :%s \n", mv.GetBody())
						if err := l.llmConsumer.Ack(context.TODO(), mv); err != nil {
							fmt.Println("ack message error: " + err.Error())
						}
					}
				}
			}
		}()
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	// wait for exit
	<-exit
	close(ch)
	wg.Wait()
	println("receiveLLMMessage close")
}

func CreateOCRConsumer(endpoint string) golang.SimpleConsumer {
	os.Setenv("mq.consoleAppender.enabled", "true")
	os.Setenv("rocketmq.client.logLevel", "error")
	golang.ResetLogger()
	simpleConsumer, err := golang.NewSimpleConsumer(&golang.Config{
		NameSpace:     "gogogo",
		ConsumerGroup: "gogogo-ocr",
		Endpoint:      endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    "",
			AccessSecret: "",
		},
	},
		golang.WithAwaitDuration(RockerMQConsumerAwaitDuration),
		golang.WithSubscriptionExpressions(map[string]*golang.FilterExpression{OCR_V1_TOPIC: golang.SUB_ALL}),
	)
	err = simpleConsumer.Start()
	if err != nil {
		logx.Error("RockerMQ-OCR-消费者创建失败")
		logx.Error(err)
		return nil
	}
	return simpleConsumer
}

func CreateLLMConsumer(endpoint string) golang.SimpleConsumer {
	os.Setenv("mq.consoleAppender.enabled", "true")
	os.Setenv("rocketmq.client.logLevel", "error")
	golang.ResetLogger()
	simpleConsumer, err := golang.NewSimpleConsumer(&golang.Config{
		NameSpace:     "gogogo",
		ConsumerGroup: "gogogo-llm",
		Endpoint:      endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    "",
			AccessSecret: "",
		},
	},
		golang.WithAwaitDuration(RockerMQConsumerAwaitDuration),
		golang.WithSubscriptionExpressions(map[string]*golang.FilterExpression{LLM_TRAIT_EXTRACT_V1_TOPIC: golang.SUB_ALL}),
	)
	err = simpleConsumer.Start()
	if err != nil {
		logx.Error("RockerMQ-LLM-消费者创建失败")
		logx.Error(err)
		return nil
	}
	return simpleConsumer
}
