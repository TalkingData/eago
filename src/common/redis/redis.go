package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const prefix = "/eago"

var Redis *RedisTool

type RedisTool struct {
	client      *redis.Client
	serviceName string
}

// 初始化连接
func InitRedis(address, password, srvName string, db int) {
	Redis = &RedisTool{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
		serviceName: srvName,
	}
}

// Set
func (rt *RedisTool) Set(key string, value interface{}, expiration time.Duration) error {
	return rt.client.Set(context.TODO(), rt.getFinalKey(key), value, expiration).Err()
}

// Del
func (rt *RedisTool) Del(key string) error {
	return rt.client.Del(context.TODO(), rt.getFinalKey(key)).Err()
}

// Expire
func (rt *RedisTool) Expire(key string, expiration time.Duration) error {
	return rt.client.PExpire(context.TODO(), rt.getFinalKey(key), expiration).Err()
}

// Get
func (rt *RedisTool) Get(key string) (string, error) {
	return rt.client.Get(context.TODO(), rt.getFinalKey(key)).Result()
}

// HasKey
func (rt *RedisTool) HasKey(key string) bool {
	if err := rt.client.Get(context.TODO(), rt.getFinalKey(key)).Err(); err == nil {
		return true
	}
	return false
}

// getFinalKey 生成KeyName
func (rt *RedisTool) getFinalKey(key string) string {
	return prefix + "/" + rt.serviceName + "/" + key
}
