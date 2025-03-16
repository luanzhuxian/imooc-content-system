package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := connRdb()
	ctx := context.Background()
	// 登录成功后，在Redis中创建一个临时的会话令牌，用于用户认证和会话管理
	// key : session_id:{user_id} val : 123  60s
	err := rdb.Set(ctx, "session_id:admin", "123", 60*time.Second).Err()
	if err != nil {
		panic(err)
	}
	sessionID, err := rdb.Get(ctx, "session_id:admin").Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	fmt.Println(sessionID)
}

func connRdb() *redis.Client {
	// redis-cli
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}
