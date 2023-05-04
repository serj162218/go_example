package initializers

import (
	"github.com/redis/go-redis/v9"
)

var (
	Redisdb *redis.Client
)

func ConnectToRedis() {
	Redisdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}