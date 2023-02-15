package redis

import (
	"github.com/go-redis/redis/v8"
)

type Option func(c *redis.Client)

// UsingOpentracingHook 使用OpentracingHook
func UsingOpentracingHook() Option {
	return func(c *redis.Client) {
		if c == nil {
			return
		}
		c.AddHook(NewOpentracingHook())
	}
}
