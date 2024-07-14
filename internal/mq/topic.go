package mq

/*
创建新topic步骤
1. 定义topic
2. 控制台手动创建topic
3. ServiceContext创建生产者, 消费者, 定义创建方法
4. gogogo.go 关闭rockerMQ生产者、开启消费者监听
5. internal/mq/producer.go 创建生产者实现
6. internal/mq/consumer.go 创建消费者实现
*/
const (
	// OCR_V1_TOPIC OCR识别topic
	OCR_V1_TOPIC = "spider-ocr-v1-topic"
	// LLM_TRAIT_EXTRACT_V1_TOPIC 特征提取topic
	LLM_TRAIT_EXTRACT_V1_TOPIC = "spider-llm-trait-extract-v1-topic"
)

const (
	OCRTag = "OCR"
	LLMTag = "LLM"
)
