package service

import (
	"log"
	"strconv"
	"strings"

	"github.com/seed30/TikTok/dao"
	"github.com/seed30/TikTok/middleware/rabbitMQ"
	"github.com/seed30/TikTok/models"

	"github.com/seed30/TikTok/middleware/redis"
	"golang.org/x/sync/singleflight"
)

var (
	gRedisFindeKeyVideo singleflight.Group
	gRedisFindKeyUser   singleflight.Group
)

type LikeServiceImpl struct {
	usi UserServiceImpl
	Vsi VideoServiceImpl
}

func (lsi *LikeServiceImpl) GetLikeByVideoIdAndUserId(videoId string, userId int64) models.Like {
	like := dao.GetLikeByVideoIdAndUserId(videoId, userId)
	return like
}

func (lsi *LikeServiceImpl) CheckLikeByVideoIdAndUserId(videoId string, userId int64) bool {
	check := dao.CheckLikeByVideoIdAndUserId(videoId, userId)
	return check
}

func (lsi *LikeServiceImpl) UpdateLike(like *models.Like) error {
	check := dao.UpdateLike(like)
	return check
}

func (lsi *LikeServiceImpl) CreateLike(like *models.Like) error {
	check := dao.CreateLike(like)
	return check
}

func (lsi *LikeServiceImpl) FindLike(like *models.Like) bool {
	exist := dao.FindLike(like)
	return exist
}

func (lsi *LikeServiceImpl) GetLikeMapByUserId(userId int64) []models.Like {
	likes := dao.GetLikeMapByUserId(userId)
	return likes
}

func (lsi *LikeServiceImpl) FindLikesByVideoId(videoId int64) ([]models.Like, error) {
	likes, err := dao.FindLikesByVideoId(videoId)
	return likes, err
}

func (lsi *LikeServiceImpl) FindLikesByUserId(userId int64) ([]models.Like, error) {
	likes, err := dao.FindLikesByUserId(userId)
	return likes, err
}

func (lsi *LikeServiceImpl) LikeAction(videoId int64, userId int64, action int) error {
	// 之前已经做过鉴权操作了， 这里就不需要再做了

	// 1 、 先查询redis 中是否存在 videoId
	// 		如果存在，那么就返回对应的value
	// 		如果不存在，则去mysql里面查询，mysql中存在的话 就将mysql中的数据 添加到redis中，不存在的话 redis创建一个新的key

	// -----------------------------------------------------
	// 这里应该需要一个锁
	// 还需要一个布隆过滤器 过滤掉 一定不存在的信息
	_, err := redis.FindKeyInVideoDB(strconv.FormatInt(videoId, 10))

	// 当stringCmd 为nil 时 必须要一个 双重锁，类似于单例模式
	// TODO 这个地方需要做并发测试
	if err == models.ErrCacheMiss {
		// 说明redis 中不存在该videoId
		// 在mysql中查询
		_, err, _ := gRedisFindeKeyVideo.Do(strconv.FormatInt(videoId, 10), func() (interface{}, error) {
			data, err := lsi.FindLikesByVideoId(videoId)
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				// 说明里面是空数据
				redis.CreateVideo(strconv.FormatInt(videoId, 10))
			} else {
				redis.SetVideo(strconv.FormatInt(videoId, 10), data)
			}

			return data, err
		})
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// -----------------------------------------------------
	_, err = redis.FindKeyInUserDB(strconv.FormatInt(userId, 10))

	// 当stringCmd 为nil 时 必须要一个 双重锁，类似于单例模式
	// TODO 这个地方需要做并发测试
	if err == models.ErrCacheMiss {
		// 说明redis 中不存在该videoId
		// 在mysql中查询
		_, err, _ := gRedisFindKeyUser.Do(strconv.FormatInt(userId, 10), func() (interface{}, error) {
			data, err := lsi.FindLikesByUserId(userId)
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				// 说明里面是空数据
				redis.CreateUser(strconv.FormatInt(userId, 10))
			} else {
				redis.SetUser(strconv.FormatInt(userId, 10), data)
			}

			return data, err
		})
		if err != nil {
			log.Println(err)
			return err
		}
	}
	// -----------------------------------------------------
	// 声明要向RabbitMQ发送的字符串
	sb := strings.Builder{}
	sb.WriteString(strconv.FormatInt(userId, 10))
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatInt(videoId, 10))
	// 在上面的操作中 已经将数据添加进了 redis 中
	// 现在就是根据 action 的操作来进行点赞。
	//首先，通过key 查询 userId 是否在集合中
	flag := redis.SIsMemberVideo(strconv.FormatInt(videoId, 10), userId)
	if flag {
		// 说明该用户已经进行了点赞的操作
		if action == 1 {
			// 此时系统出现了问题 ， 如果用户已经点过赞了，然后还发出action为1的操作 就需要报错
			log.Println("用户已经点赞了该视频，不能重复点赞")
			return models.ErrLikeAction
		} else if action == 2 {
			// 此时用户取消点赞，那么redis 需要在set中删除该用户的信息
			err := redis.SRemVideo(strconv.FormatInt(videoId, 10), userId)
			if err != nil {
				return err
			}
			// 此时已经删除了redis中的 视频对应的 该用户点赞信息
			// 还需要对 RedisUserLikeDb 进行操作
			err = redis.SRemUser(strconv.FormatInt(userId, 10), videoId)
			if err != nil {
				err1 := redis.SAddUserToVideo(strconv.FormatInt(videoId, 10), userId)
				if err1 != nil {
					return models.ErrVideoUserConsist
				}
				return err
			} else {
				// 向RmqLikeDel 发送请求
				// 数据库同步取消点赞
				rabbitMQ.RmqLikeDel.Publish(sb.String())
			}
		}
	} else {
		// flag 为 false 的情况时  , 如果action 为1 ，那么就是点赞
		// action 为 2 时 ，就要 返回一个服务错误
		if action == 1 {
			// 此时用户开始点赞
			// 想 RedisVideoLikeDb 中的对应key 的set中添加value
			err := redis.SAddUserToVideo(strconv.FormatInt(videoId, 10), userId)
			if err != nil {
				return err
			}
			err = redis.SAddVideoToUser(strconv.FormatInt(userId, 10), videoId)
			if err != nil {
				// 就需要对 RedisVideoLikeDb 进行回滚操作
				err1 := redis.SRemVideo(strconv.FormatInt(videoId, 10), userId)
				if err1 != nil {
					return err1
				}
				return err
			} else {
				// 向RmqLikeDel 发送请求
				// 数据库同步取消点赞
				rabbitMQ.RmqLikeAdd.Publish(sb.String())
			}
		}
	}

	return nil
}

// FindFavoriteVideos 获取喜欢列表
func (lsi *LikeServiceImpl) FindFavoriteVideos(userId int64) ([]models.Video, error) {
	// 首先从Redis里面去查询
	// 首先在Redis中查询是否存在这个key
	exist, _ := redis.ExistKeyInUserDB(strconv.FormatInt(userId, 10))
	if exist == 1 {
		// 说明数据存在 ， 那么就需要获取Redis里面的所有key所对应的value值
		valueStrings := redis.FindSetInUserDB(strconv.FormatInt(userId, 10))
		valueIds := make([]int64, len(valueStrings))
		for i := 1; i < len(valueStrings); i++ {
			// 为了使得Redis 中的key不会被删除，在第一行中添加了一个 -1 的value， 需要把这个-1值给排除在外
			valueIds[i], _ = strconv.ParseInt(valueStrings[i], 10, 64)
		}
		videos, err := lsi.Vsi.GetVideosByVideoIds(valueIds)
		if err != nil {
			log.Println("从数据库中获取点赞的视频列表失败")
			return nil, err
		}
		return videos, nil
	} else {
		// 说明数据不存在
		// 如果数据不存在的话 ， 那么就很有必要向数据库进行查询，并且将查询的结果保存至Redis中。
		// 先向数据库进行查询

		likes, err := lsi.FindLikesByUserId(userId)
		redis.SetUser(strconv.FormatInt(userId, 10), likes)
		if err != nil {
			log.Println("从数据库likes表中获取信息失败")
			return nil, err
		}
		videoIds := make([]int64, len(likes))
		for i := 1; i < len(likes); i++ {
			videoIds[i] = likes[i].VideoId
		}
		videos, err := lsi.Vsi.GetVideosByVideoIds(videoIds)
		if err != nil {
			log.Println("从数据库likes表中获取信息失败")
			return nil, nil
		}
		return videos, nil
	}
}
