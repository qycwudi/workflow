package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ DataModel = (*customDataModel)(nil)

type (
	// DataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDataModel.
	DataModel interface {
		dataModel
		InsertOne(ctx context.Context, data *Data) error
		UpdateOcrResultByKey(ctx context.Context, data *Data) (*mongo.UpdateResult, error)
		UpdateLlmResultByKey(ctx context.Context, data *Data) (*mongo.UpdateResult, error)
	}

	customDataModel struct {
		*defaultDataModel
	}
)

func (m customDataModel) UpdateOcrResultByKey(ctx context.Context, data *Data) (*mongo.UpdateResult, error) {
	// 创建filter，用于查找要更新的文档
	filter := bson.M{"key": data.Key}

	// 创建update，用于指定要更新的字段
	update := bson.M{
		"$set": bson.M{
			"ocrResult": data.OcrResult,
			"updateAt":  time.Now(),
		},
	}
	return m.conn.UpdateOne(ctx, filter, update)
}

func (m customDataModel) UpdateLlmResultByKey(ctx context.Context, data *Data) (*mongo.UpdateResult, error) {
	// 创建filter，用于查找要更新的文档
	filter := bson.M{"key": data.Key}

	// 创建update，用于指定要更新的字段
	update := bson.M{
		"$set": bson.M{
			"llmResult": data.OcrResult,
			"updateAt":  time.Now(),
		},
	}
	return m.conn.UpdateOne(ctx, filter, update)
}

func (m customDataModel) InsertOne(ctx context.Context, data *Data) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}
	_, err := m.conn.InsertOne(ctx, data)
	return err
}

// NewDataModel returns a model for the mongo.
func NewDataModel(url string) DataModel {
	db := "spider"
	collection := "data"
	conn := mon.MustNewModel(url, db, collection, mon.WithTimeout(2*time.Second))
	return &customDataModel{
		defaultDataModel: newDefaultDataModel(conn),
	}
}
