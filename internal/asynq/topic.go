package asynq

// A list of task types.
const (
	// TypeOCRRecognize OCR识别
	TypeOCRRecognize = "ocr:recognize:text"

	// TypeLLMFeatureExtraction 大模型提取特征
	TypeLLMFeatureExtraction = "llm:feature:extraction"
)
