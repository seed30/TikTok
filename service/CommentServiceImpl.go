package service

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/seed30/TikTok/dao"
	"github.com/seed30/TikTok/middleware/rabbitMQ"
	"github.com/seed30/TikTok/middleware/redis"
	"github.com/seed30/TikTok/models"
)

type CommentServiceImpl struct {
}

func (csi *CommentServiceImpl) SaveComment(comment *models.Comment) error {
	err := dao.SaveComment(comment)
	return err
}

func (csi *CommentServiceImpl) FindCommentByCommentId(commentId int64) (models.Comment, error) {
	comment, err := dao.FindCommentByCommentId(commentId)

	return comment, err
}

func (csi *CommentServiceImpl) DeletComment(comment *models.Comment) error {
	err := dao.DeletComment(comment)

	return err
}

func (csi *CommentServiceImpl) FindAllCommentByVideoId(videoId int64) ([]models.Comment, error) {
	var comments []models.Comment
	// 首先向Redis里面查询
	_, err := redis.FindKeyInCommentDB(strconv.FormatInt(videoId, 10))
	if err != nil && err == models.ErrCacheMiss {
		// 说明数据中不存在 当前视频Id 的key
		// 向 mysql 中查询
		_, err, _ := gRedisFindeKeyVideo.Do(strconv.FormatInt(videoId, 10), func() (interface{}, error) {
			data, err := csi.FindAllCommentByVideoId(videoId)
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				// 说明里面是空数据 ， 也就是当前视频从来没有人点赞
				redis.CreateComment(strconv.FormatInt(videoId, 10))
			} else {
				// 将评论id放入redis中
				redis.SetComments(strconv.FormatInt(videoId, 10), data)
			}

			return data, err
		})
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	// redis 中存在
	data, err := redis.HGetAll(strconv.FormatInt(videoId, 10))
	if err != nil {
		panic(err)
	}
	// data是一个map类型，这里使用使用循环迭代输出

	for _, val := range data {
		var comment models.Comment
		err := json.Unmarshal([]byte(val), &comment)
		if err != nil {
			return nil, nil
		}
		comments = append(comments, comment)

	}

	return comments, err
}

func (csi *CommentServiceImpl) AddComment(videoId int64, comment models.Comment) error {
	// 这是添加评论的操作
	// 首先向Redis中查询评论数据，与点赞功能的实现类似
	// 不过需要向Redis中存储的不再是用户ID，而是评论的整个结构体，方便之后CommentList操作时，全部取出
	_, err := redis.FindKeyInCommentDB(strconv.FormatInt(videoId, 10))
	if err != nil && err == models.ErrCacheMiss {
		// 说明数据中不存在 当前视频Id 的key
		// 向 mysql 中查询
		_, err, _ := gRedisFindeKeyVideo.Do(strconv.FormatInt(videoId, 10), func() (interface{}, error) {
			data, err := csi.FindAllCommentByVideoId(videoId)
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				// 说明里面是空数据 ， 也就是当前视频从来没有人点赞
				redis.CreateComment(strconv.FormatInt(videoId, 10))
			} else {
				// 将评论id放入redis中
				redis.SetComments(strconv.FormatInt(videoId, 10), data)
			}

			return data, err
		})
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// 必须先操作数据库才能够拿到commentId
	err = csi.SaveComment(&comment)
	if err != nil {
		return err
	}
	// 准备之后需要想RabbitMQ发送的请求
	commentJSON, _ := json.Marshal(comment)
	log.Println("commentJSON : ", commentJSON)

	err = redis.HSetComment(strconv.FormatInt(videoId, 10), strconv.FormatInt(comment.Id, 10), commentJSON)
	if err != nil {
		// 如果有报错，那么就不用想数据库发送请求
		return errors.New("评论id添加到redis中失败")
	}

	//err = csi.SaveComment(&comment)
	//if err != nil {
	//	return errors.New("评论添加到数据库中失败")
	//}

	//rabbitMQ.RmqCommentAdd.Publish(string(commentJSON))

	// TODO 这里要考虑数据一致性的问题
	return nil
}

func (csi *CommentServiceImpl) DelComment(comment models.Comment) error {
	videoId := comment.VideoId
	_, err := redis.FindKeyInCommentDB(strconv.FormatInt(videoId, 10))
	if err != nil && err == models.ErrCacheMiss {
		// 说明数据中不存在 当前视频Id 的key
		// 向 mysql 中查询
		_, err, _ := gRedisFindeKeyVideo.Do(strconv.FormatInt(videoId, 10), func() (interface{}, error) {
			data, err := csi.FindAllCommentByVideoId(videoId)
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				// 说明里面是空数据 ， 也就是当前视频从来没有人点赞
				redis.CreateComment(strconv.FormatInt(videoId, 10))
				return nil, errors.New("comment is not exist")
			} else {
				redis.SetComments(strconv.FormatInt(videoId, 10), data)
			}

			return data, err
		})
		if err != nil {
			log.Println(err)
			return err
		}
	}

	commentId := comment.Id

	err = redis.HDelComment(strconv.FormatInt(videoId, 10), strconv.FormatInt(commentId, 10))
	if err != nil {
		return err
	}

	rabbitMQ.RmqCommentDel.Publish(strconv.FormatInt(commentId, 10))

	return nil
}
