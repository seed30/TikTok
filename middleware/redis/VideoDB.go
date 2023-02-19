package redis

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/seed30/TikTok/models"
)

func FindKeyInVideoDB(videoId string) (*redis.StringCmd, error) {
	stringCmd := RedisVideoLikeDb.Get(videoId)

	if stringCmd.Err() == redis.Nil {
		return stringCmd, models.ErrCacheMiss
	}

	return stringCmd, nil
}

func SetVideo(videoId string, data []models.Like) {
	pipe := RedisVideoLikeDb.Pipeline()
	pipe.SAdd(videoId, -1)
	for i := 0; i < len(data); i++ {
		userId := data[i].UserId
		pipe.SAdd(videoId, userId)
	}
	_, err := pipe.Exec()
	if err != nil {
		log.Println("save the data into redis failed !")
	}
}

func CreateVideo(videoId string) {
	err := RedisVideoLikeDb.SAdd(videoId, -1).Err()
	if err != nil {
		log.Println(" create new key failed in redis ")
	}
}

func SIsMemberVideo(videoId string, userId int64) bool {
	flag, err := RedisVideoLikeDb.SIsMember(videoId, userId).Result()
	if flag {
		log.Println("redis 中存在当前用户 ，且当前用户已经点赞该条视频")
	} else {
		log.Println("redis 中不存在当前用户， 说明该用户没有给该视频点赞")
	}
	if err != nil {
		log.Println("redis 查询key出错，请检查redis")
	}
	return flag
}

func SRemVideo(videoId string, userId int64) error {
	_, err := RedisVideoLikeDb.SRem(videoId, userId).Result()
	if err != nil {
		return err
	}
	return nil
}

func SAddUserToVideo(videoId string, userId int64) error {
	err := RedisVideoLikeDb.SAdd(videoId, userId).Err()
	return err
}
