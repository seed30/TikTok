package redis

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/seed30/TikTok/config"
)

var RedisVideoLikeDb *redis.Client

var RedisUserLikeDb *redis.Client

var RedisCommentDb *redis.Client

func InitClient() (err error) {
	// 视频ID  对应 多个点赞人ID
	RedisVideoLikeDb = redis.NewClient(&redis.Options{
		Addr:         config.RedisAddr,     // redis地址
		Password:     config.RedisPassword, // redis密码，没有则留空
		DB:           0,                    // 默认数据库，默认是0
		PoolSize:     config.RedisPoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	pong, err := RedisVideoLikeDb.Ping().Result()
	if err != nil {
		return err
	}
	log.Println("redis init 0  : ", pong)

	// 用户ID 对应 多个已点赞的视频ID
	RedisUserLikeDb = redis.NewClient(&redis.Options{
		Addr:         config.RedisAddr,     // redis地址
		Password:     config.RedisPassword, // redis密码，没有则留空
		DB:           1,                    // 1号数据库，默认是0
		PoolSize:     config.RedisPoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	pong, err = RedisUserLikeDb.Ping().Result()
	if err != nil {
		return err
	}
	log.Println("redis init 1 : ", pong)

	// 视频ID 对应 多个评论信息
	RedisCommentDb = redis.NewClient(&redis.Options{
		Addr:         config.RedisAddr,
		Password:     config.RedisPassword,
		DB:           2,
		PoolSize:     config.RedisPoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	return nil

}
