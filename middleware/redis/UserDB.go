package redis

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/seed30/TikTok/models"
)

func FindKeyInUserDB(userId string) (*redis.StringCmd, error) {
	stringCmd := RedisUserLikeDb.Get(userId)
	if stringCmd.Err() == redis.Nil {
		return stringCmd, models.ErrCacheMiss
	}
	return stringCmd, nil
}

func ExistKeyInUserDB(userId string) (int64, error) {
	exists, err := RedisUserLikeDb.Exists(userId).Result()
	return exists, err
}

func FindSetInUserDB(userId string) []string {
	result, err := RedisUserLikeDb.SMembers(userId).Result()
	if err != nil {
		log.Println(" Redis Find set has something wrong")
	}
	return result
}

func CreateUser(userId string) {
	err := RedisUserLikeDb.SAdd(userId, -1).Err()
	if err != nil {
		log.Println(" create new key failed in redis ")
	}
}

func SetUser(userId string, data []models.Like) {

	// 使用管道 批量添加数据 会更加快速
	pipe := RedisUserLikeDb.Pipeline()
	pipe.SAdd(userId, -1)
	for i := 0; i < len(data); i++ {
		videoId := data[i].VideoId
		pipe.SAdd(userId, videoId)
	}
	_, err := pipe.Exec()

	if err != nil {
		log.Println("save the data into redis failed !")
		panic(err)
	}
}
func SRemUser(userId string, videoId int64) error {
	_, err := RedisUserLikeDb.SRem(userId, videoId).Result()

	if err != nil {
		return err
	}
	return nil
}

func SAddVideoToUser(userId string, videoId int64) error {
	err := RedisUserLikeDb.SAdd(userId, videoId).Err()
	return err
}
