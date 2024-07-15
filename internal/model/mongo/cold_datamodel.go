package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ ColdDataModel = (*customColdDataModel)(nil)

type (
	// ColdDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customColdDataModel.
	ColdDataModel interface {
		coldDataModel
		InsertOne(ctx context.Context, data *ColdData) error
		UpdateOcrResultByKey(ctx context.Context, data *ColdData) (*mongo.UpdateResult, error)
		UpdateLlmResultByKey(ctx context.Context, data *ColdData) (*mongo.UpdateResult, error)
		FindOneByKey(ctx context.Context, key string) (*ColdData, error)
	}

	customColdDataModel struct {
		*defaultColdDataModel
	}
)

func (m customColdDataModel) FindOneByKey(ctx context.Context, key string) (*ColdData, error) {

	var data ColdData
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

func (m customColdDataModel) UpdateOcrResultByKey(ctx context.Context, data *ColdData) (*mongo.UpdateResult, error) {
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

func (m customColdDataModel) UpdateLlmResultByKey(ctx context.Context, data *ColdData) (*mongo.UpdateResult, error) {
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

func (m customColdDataModel) InsertOne(ctx context.Context, data *ColdData) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}
	_, err := m.conn.InsertOne(ctx, data)
	return err
}

// NewColdDataModel returns a model for the mongo.
func NewColdDataModel(url string) ColdDataModel {
	db := "spider"
	collection := "cold-data"
	conn := mon.MustNewModel(url, db, collection, mon.WithTimeout(2*time.Second))
	return &customColdDataModel{
		defaultColdDataModel: newDefaultColdColdDataModel(conn),
	}
}
