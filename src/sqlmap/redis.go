package sqlmap

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var (
	RedisClient2 *redis.Client
	RedisClient  *redis.Client
	Ctx          context.Context
)

func init() {
	log.Println("初始化Redis")
	Ctx = context.Background()
	RedisCli()
	RedisCli2()
}

func RedisCli() {
	RedisClient = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     rdb.Addr,
		Password: rdb.Password,
		DB:       rdb.Db1,

		PoolSize:     15,
		MinIdleConns: 10,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,

		IdleCheckFrequency: 60 * time.Second,
		IdleTimeout:        5 * time.Second,
		MaxConnAge:         0 * time.Second,

		MaxRetries:      0,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})
	err = RedisClient.Ping(Ctx).Err()
	if err != nil {
		log.Println(err.Error())
	}
}

func RedisCli2() {
	RedisClient2 = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     rdb.Addr,
		Password: rdb.Password,
		DB:       rdb.Db2,

		PoolSize:     15,
		MinIdleConns: 10,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,

		IdleCheckFrequency: 60 * time.Second,
		IdleTimeout:        5 * time.Second,
		MaxConnAge:         0 * time.Second,

		MaxRetries:      0,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})
	err = RedisClient2.Ping(Ctx).Err()
	if err != nil {
		log.Println(err.Error())
	}
}
