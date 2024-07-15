package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ HotDataModel = (*customHotHotDataModel)(nil)

type (
	// HotDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHotHotDataModel.
	HotDataModel interface {
		hotDataModel
		InsertOne(ctx context.Context, data *HotData) error
		UpdateOcrResultByKey(ctx context.Context, data *HotData) (*mongo.UpdateResult, error)
		UpdateLlmResultByKey(ctx context.Context, data *HotData) (*mongo.UpdateResult, error)
		FindOneByKey(ctx context.Context, key string) (*HotData, error)
	}

	customHotHotDataModel struct {
		*defaultHotHotDataModel
	}
)

func (m customHotHotDataModel) FindOneByKey(ctx context.Context, key string) (*HotData, error) {

	var data HotData
	err := m.conn.FindOne(ctx, &data, bson.M{"key": key})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m customHotHotDataModel) UpdateOcrResultByKey(ctx context.Context, data *HotData) (*mongo.UpdateResult, error) {
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

func (m customHotHotDataModel) UpdateLlmResultByKey(ctx context.Context, data *HotData) (*mongo.UpdateResult, error) {
	// 创建filter，用于查找要更新的文档
	filter := bson.M{"key": data.Key}

	// 创建update，用于指定要更新的字段
	update := bson.M{
		"$set": bson.M{
			"llmResult": data.LLMResult,
			"updateAt":  time.Now(),
		},
	}
	return m.conn.UpdateOne(ctx, filter, update)
}

func (m customHotHotDataModel) InsertOne(ctx context.Context, data *HotData) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}
	_, err := m.conn.InsertOne(ctx, data)
	return err
}

// NewHotDataModel returns a model for the mongo.
func NewHotDataModel(url string) HotDataModel {
	db := "spider"
	collection := "hot-data"
	conn := mon.MustNewModel(url, db, collection, mon.WithTimeout(2*time.Second))
	return &customHotHotDataModel{
		defaultHotHotDataModel: newDefaultHotHotDataModel(conn),
	}
}
