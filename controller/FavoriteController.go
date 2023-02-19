package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seed30/TikTok/models"
	"github.com/seed30/TikTok/service"
)

/*
*
点赞系统需要考虑：

	1、Redis 采用 bitmap 还是 set来存储 点赞信息。 -----》 set适合小访问量 ， bitmap 适合百万级别的访问量（经验值）
	2、点赞需要考虑两种信息：
		2.1 视频的点赞信息，有哪些用户点赞了该视频
		2.2 用户的点赞信息，该用户点赞了哪些视频
		基于以上两种信息，就要考虑是分库还是两种信息混合存储，除此之外 还要考虑事务的一致性
	3、 Redis持久化方法。
*/
func FavoriteAction(c *gin.Context) {
	/*
		要处理的几件事情：
			1、点击后，数据库中的 喜欢数据要发生改变，改变如下：
				如果喜欢数据存在，则修改 喜欢列表中的 喜欢操作， 喜欢为1 不喜欢为0
				如果喜欢数据不存在， 则创建喜欢数据
			2、以上操作要同步到 video 数据中，并操作 favorite_count操作
	*/

	// TODO 这个操作肯定要放到Redis 里面的，如果不放入的话 会频繁的访问数据库，导致数据库压力过大

	// TODO 鉴权操作 省略

	token := c.Query("token")
	actionString := c.Query("action_type")
	action, _ := strconv.Atoi(actionString)

	log.Println("action:", action)
	videoIdString := c.Query("video_id")

	videoId, _ := strconv.ParseInt(videoIdString, 10, 64)

	// 获取用户ID
	usi := service.UserServiceImpl{}
	_, user := usi.GetTableUserByToken(token)
	userId := user.Id

	lsi := service.LikeServiceImpl{}
	//vsi := service.VideoServiceImpl{}

	err := lsi.LikeAction(videoId, userId, action)

	if err == models.ErrLikeAction {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: models.ErrLikeActionCode,
			StatusMsg:  "the like action is wrong",
		})
	}

	// 根据 action 的值来判断是否需要进行操作
	//like := models.Like{
	//	UserId:  userId,
	//	VideoId: int64(videoId),
	//	Cancel:  action,
	//}
	//exist := lsi.FindLike(&like)
	//log.Println(like)
	//if !exist && action != 2 {
	//	// 数据不存在 并且表示喜欢
	//	log.Println("No like data")
	//	err := lsi.CreateLike(&like)
	//	if err != nil {
	//		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "Create like failed"})
	//		return
	//	}
	//} else {
	//	like.Cancel = action
	//	err := lsi.UpdateLike(&like)
	//	if err != nil {
	//		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "Update like failed"})
	//		return
	//	}
	//}

	// TODO 这个操作可以进行优化
	//vsi.UpdateFavoriteByVideoId(int64(videoId))

	c.JSON(http.StatusOK, models.Response{StatusCode: 0})
}

func FavoriteList(c *gin.Context) {
	// TODO  鉴权
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	lsi := service.LikeServiceImpl{}

	favoriteVideos, err := lsi.FindFavoriteVideos(userId)
	if err != nil {
		c.JSON(http.StatusOK, models.VideoListResponse{
			Response: models.Response{
				StatusCode: models.ErrReids,
			},
		})
		return
	}

	//videos := lsi.Vsi.GetAllVideos()
	//likes := lsi.GetLikeMapByUserId(userId)
	//
	//var set map[int64]struct{}
	//set = make(map[int64]struct{})
	//
	//for _, value := range likes {
	//	set[value.VideoId] = struct{}{}
	//}
	//favoriteVideos := make([]models.Video, len(likes))
	//index := 0
	//for i := 0; i < len(videos); i++ {
	//	if _, ok := set[videos[i].Id]; ok {
	//		favoriteVideos[index] = videos[i]
	//		index++
	//	}
	//}
	//log.Println(favoriteVideos)
	c.JSON(http.StatusOK, models.VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: favoriteVideos,
	})

}
