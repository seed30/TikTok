package dao

import (
	"errors"

	"github.com/seed30/TikTok/models"
	"gorm.io/gorm"
)

func SaveComment(comment *models.Comment) error {
	//err := Db.Table("comments").Save(&comment).Error
	//return err

	err := Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Table("comments").Save(&comment).Error; err != nil {
			return err
		}
		var comments []models.Comment
		if err := tx.Table("comments").Where("video_id = ?", comment.VideoId).Where("cancel = ?", 1).Find(&comments).Error; err != nil {
			return err
		}
		count := len(comments)

		if err := tx.Table("videos").Where("id = ?", comment.VideoId).Update("comment_count", count).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func FindCommentByCommentId(commentId int64) (models.Comment, error) {
	var comment models.Comment
	result := Db.Table("comments").Preload("User").First(&comment, "id = ?", commentId)
	if result.RowsAffected > 0 {
		return comment, nil
	} else {
		return comment, errors.New("评论不存在")
	}
}

func DeletComment(comment *models.Comment) error {

	err := Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Model(&comment).Table("comments").Where("id = ?", comment.Id).Update("cancel", 2).Error; err != nil {
			return err
		}
		var comments []models.Comment
		if err := tx.Table("comments").Where("video_id = ?", comment.VideoId).Where("cancel = ?", 1).Find(&comments).Error; err != nil {
			return err
		}
		count := len(comments)

		if err := tx.Table("videos").Where("id = ?", comment.VideoId).Update("comment_count", count).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func FindAllCommentByVideoId(comments *[]models.Comment, videoId int64) error {
	err := Db.Model(&models.Comment{}).Preload("User").Where("video_id = ?", videoId).Where("cancel = ?", 1).Find(&comments).Error
	return err
}
