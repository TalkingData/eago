package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const prefix = "/td/eago"

var Redis *RedisTool

type RedisTool struct {
	client      *redis.Client
	serviceName string
}

// InitRedis 初始化Redis
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

// Close 关闭Redis
func Close() {
	if Redis == nil {
		return
	}
	_ = Redis.client.Close()
}

func (rt *RedisTool) Set(key string, value interface{}, expiration time.Duration) error {
	return rt.client.Set(context.Background(), rt.getFinalKey(key), value, expiration).Err()
}

func (rt *RedisTool) Del(key string) error {
	return rt.client.Del(context.Background(), rt.getFinalKey(key)).Err()
}

func (rt *RedisTool) Expire(key string, expiration time.Duration) error {
	return rt.client.PExpire(context.Background(), rt.getFinalKey(key), expiration).Err()
}

func (rt *RedisTool) Get(key string) (string, error) {
	return rt.client.Get(context.Background(), rt.getFinalKey(key)).Result()
}

func (rt *RedisTool) HasKey(key string) bool {
	if err := rt.client.Get(context.Background(), rt.getFinalKey(key)).Err(); err == nil {
		return true
	}
	return false
}

// getFinalKey 生成KeyName
func (rt *RedisTool) getFinalKey(key string) string {
	return prefix + "/" + rt.serviceName + "/" + key
}
