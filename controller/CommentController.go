package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seed30/TikTok/models"
	"github.com/seed30/TikTok/service"
	"github.com/seed30/TikTok/utils"
)

func CommentAction(c *gin.Context) {
	// 评论操作
	// TODO 鉴权操作
	token := c.Query("token")

	usi := service.UserServiceImpl{}
	_, user := usi.GetTableUserByToken(token)

	// 获取 POST 信息
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	action, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	commentText := c.Query("comment_text")
	commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64) // 要删除评论的id 在action=2的时候使用

	csi := service.CommentServiceImpl{}
	var comment models.Comment
	// 判断 action ， action = 1 -> 评论操作 ；  action = 2 -> 删除评论操作(只有发布者和视频拥有者可以删除)

	if action == 1 {
		// 过滤敏感词
		hasSensitive, SensitiveText := utils.Filter.Validate(commentText)
		if hasSensitive == false {
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  "(" + SensitiveText + ") 为敏感词，发表评论失败",
			})
			return
		}

		comment = models.Comment{
			User:        user,
			UserId:      user.Id,
			VideoId:     videoId,
			Cancel:      1,
			Content:     commentText,
			CreatedDate: time.Now(),
		}
		// 添加评论的操作
		err := csi.AddComment(videoId, comment)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  "评论失败",
			})

			return
		}
	} else if action == 2 {

		comment = models.Comment{
			Id:      commentId,
			User:    user,
			UserId:  user.Id,
			VideoId: videoId,
			Cancel:  2,
			Content: commentText,
		}

		err := csi.DelComment(comment)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  "评论失败",
			})
			return
		}

		// 删除评论的操作
	}

	c.JSON(http.StatusOK, models.CommentActionResponse{
		Response: models.Response{StatusCode: 0},
		Comment:  comment,
	})
}

func CommentList(c *gin.Context) {
	// 展示评论列表

	// TODO 鉴权

	// 根据 video_id 查询所有相关的评论
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	csi := service.CommentServiceImpl{}
	comments, err := csi.FindAllCommentByVideoId(videoId)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  "获取评论列表失败，请稍等",
		})
	}
	c.JSON(http.StatusOK, models.CommentListResponse{
		Response:    models.Response{StatusCode: 0},
		CommentList: comments,
	})
}
