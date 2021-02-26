package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const fixedPrefix = "eago"

var (
	Redis RedisTool

	serviceName string
	ctx         = context.Background()
)

type RedisTool struct {
	Client *redis.Client
}

// 初始化连接
func InitRedis(address string, password string, db int, srvName string) error {
	Redis.Client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	serviceName = srvName
	return nil
}

func (rt *RedisTool) Set(key string, value interface{}, expiration time.Duration) error {
	return rt.Client.Set(ctx, rt.getFinalKey(key), value, expiration).Err()
}

func (rt *RedisTool) Del(key string) error {
	return rt.Client.Del(ctx, rt.getFinalKey(key)).Err()
}

func (rt *RedisTool) Expire(key string, expiration time.Duration) error {
	return rt.Client.PExpire(ctx, rt.getFinalKey(key), expiration).Err()
}

func (rt *RedisTool) Get(key string) (string, error) {
	return rt.Client.Get(ctx, rt.getFinalKey(key)).Result()
}

func (rt *RedisTool) HasKey(key string) bool {
	if err := rt.Client.Get(ctx, rt.getFinalKey(key)).Err(); err == nil {
		return true
	}
	return false
}

// 生成
func (rt RedisTool) getFinalKey(key string) string {
	return fixedPrefix + "/" + serviceName + "/" + key
}
