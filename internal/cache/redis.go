package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	ApiPrefixRedisKey = "api:%s:"
	ApiInfoRedisKey   = ApiPrefixRedisKey + "info"

	ApiSecretKeyPrefixRedisKey = ApiPrefixRedisKey + "secret_key:"
	ApiSecretKeyRedisKey       = ApiSecretKeyPrefixRedisKey + "%s"
)

type RedisCache struct {
	client redis.UniversalClient
}

var Redis *RedisCache

func NewRedis(c redis.UniversalClient) {
	Redis = &RedisCache{
		client: c,
	}
}

// Set 设置缓存,带过期时间
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Subscribe 订阅频道
func (r *RedisCache) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.client.Subscribe(ctx, channel)
}

// Publish 发布消息到频道
func (r *RedisCache) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.client.Publish(ctx, channel, message).Err()
}

// Del 删除缓存
// 如果key不存在,Del方法会返回0,但不会返回错误
func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DelByPrefix 根据前缀删除缓存
// 如果key不存在,Del方法会返回0,但不会返回错误
func (r *RedisCache) DelByPrefix(ctx context.Context, prefix string) error {
	// 使用 SCAN 命令迭代所有匹配的 key
	iter := r.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}

	// 如果有匹配的key,批量删除
	if len(keys) > 0 {
		logx.Infof("delete redis keys: %v", keys)
		return r.client.Del(ctx, keys...).Err()
	}
	return nil
}
