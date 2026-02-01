package redis

import (
	"github.com/go-redis/redis/v8"
	"os"
)

func InitRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379" // fallback для локалки
	}

	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
}
