package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotData struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UpdateAt  time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt  time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
	Key       string             `bson:"key" json:"key"`
	Source    string             `bson:"source" json:"source"`
	OcrResult string             `bson:"ocrResult" json:"ocrResult"`
	LLMResult map[string]string  `bson:"llmResult" json:"LLMResult"`
}
