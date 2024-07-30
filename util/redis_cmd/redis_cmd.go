package redis_cmd

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

func NewRedisCmd(addr string) redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		ReadTimeout:  time.Millisecond * 200,
		WriteTimeout: time.Millisecond * 200,
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		log.Fatalf("redis connect error: %v", err)
	}
	return client
}
