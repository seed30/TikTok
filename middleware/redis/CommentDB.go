package redis

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/seed30/TikTok/models"
)

func FindKeyInCommentDB(videoId string) (*redis.StringCmd, error) {
	StringCmd := RedisCommentDb.Get(videoId)

	if StringCmd.Err() == redis.Nil {
		return StringCmd, models.ErrCacheMiss
	}
	return StringCmd, nil
}

func CreateComment(videoId string) {

	// 将这个空的Comment 保存到Redis中
	err := RedisCommentDb.SAdd(videoId, -1).Err()

	if err != nil {
		log.Println(" create new key in redis failed !")
	}

}

func SetComments(videoId string, comments []models.Comment) {
	// 将批量的数据添加到Redis中
	pipe := RedisCommentDb.Pipeline()

	pipe.HSet(videoId, "0", "nil")

	for i := 0; i < len(comments); i++ {
		commentJSON, _ := json.Marshal(comments[i])
		pipe.HSet(videoId, strconv.FormatInt(comments[i].Id, 10), commentJSON)
	}

	_, err := pipe.Exec()
	if err != nil {
		log.Println("Set Comments into Redis failed ! ")
		return
	}
}

func HSetComment(videoId string, commentId string, commentJSON []byte) error {
	err := RedisCommentDb.HSet(videoId, commentId, commentJSON).Err()
	return err
}

func HDelComment(videoId string, commentId string) error {
	err := RedisCommentDb.HDel(videoId, commentId).Err()
	return err
}

func HGetAll(videoId string) (map[string]string, error) {
	data, err := RedisCommentDb.HGetAll(videoId).Result()
	return data, err
}

func SAddComment(videoId string, commentId int64) error {
	err := RedisCommentDb.SAdd(videoId, commentId).Err()
	if err != nil {
		log.Println("Add Comment into Redis failed ! ")
		return err
	}
	return nil
}
func SRemComment(videoId string, commentId int64) error {
	err := RedisCommentDb.SRem(videoId, commentId).Err()

	if err != nil {
		log.Println(" Delet Comment in Redis failed !")
	}
	return nil
}
