package cache

import (
	"context"
	"fmt"
	"time"

	config "main/config"

	"github.com/go-redis/redis/v8"
)

const Nil = redis.Nil // empty val

var (
	rdb *redis.Client
	ctx context.Context
)

// Set 设置缓存字段
func Set(key, value string, ttl time.Duration) error {
	return rdb.Set(ctx, key, value, ttl*time.Second).Err()
}

// Get 获取缓存字段
func Get(key string) (val string, err error) {
	val, err = rdb.Get(ctx, key).Result()
	return
}

// Del 删除缓存字段
func Del(key string) error {
	return rdb.Del(ctx, key).Err()
}

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Address,
		Password: config.Conf.Redis.Password,
		Network:  config.Conf.Redis.Network,
		DB:       config.Conf.Redis.DB,
	})
	ctx = context.Background()
	pong := rdb.Ping(ctx)
	if err := pong.Err(); err != nil {
		panic(err)
	} else {
		fmt.Printf("%s, connect to redis[%s]\n", pong.Val(), config.Conf.Redis.Address)
	}
}
