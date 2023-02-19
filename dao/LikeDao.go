package dao

import (
	"log"

	"github.com/seed30/TikTok/models"
	"gorm.io/gorm"
)

func GetLikeByVideoIdAndUserId(videoId string, userId int64) models.Like {
	var like models.Like
	Db.Table("likes").First(&like, "video_id = ?", videoId, "user_id = ?", userId)
	return like
}

func CheckLikeByVideoIdAndUserId(videoId string, userId int64) bool {
	var like models.Like
	first := Db.Table("likes").First(&like, "video_id = ?", videoId, "user_id = ?", userId)
	if first.Error != nil {
		return false
	} else {
		return true
	}
}

func UpdateLike(like *models.Like) error {

	err := Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&like).Where("id = ?", like.Id).Update("Cancel", like.Cancel).Error; err != nil {
			log.Println(err)
			return err
		}
		// 更新videos表
		result := tx.Table("likes").Where("video_id = ?", like.VideoId).Where("cancel = ?", 1).Find(&[]models.Like{})
		if result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}
		count := result.RowsAffected

		if err := tx.Table("videos").Where("id = ?", like.VideoId).Update("favorite_count", count).Error; err != nil {
			log.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateLike(like *models.Like) error {
	result := Db.Table("likes").Create(&like)
	return result.Error
}

func FindLike(like *models.Like) bool {
	find := Db.Table("likes").Where("video_id = ?", like.VideoId).Where("user_id = ?", like.UserId).Find(&like)
	if find.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func SumFavoriteCountByVideoId(videoId int64) int64 {
	result := Db.Table("likes").Where("video_id = ?", videoId).Where("cancel = ?", 1).Find(&[]models.Like{})
	return result.RowsAffected
}

func GetLikeMapByUserId(userId int64) []models.Like {
	likes := []models.Like{}
	Db.Table("likes").Where("user_id = ?", userId).Where("cancel = ?", 1).Find(&likes)
	return likes
}
func SumCommentCountByVideoId(videoId int64) int64 {
	result := Db.Table("comments").Where("video_id = ?", videoId).Where(&models.Comment{Cancel: 0}).Find(&[]models.Comment{})
	return result.RowsAffected
}

func FindLikesByVideoId(videoId int64) ([]models.Like, error) {
	var likes []models.Like
	result := Db.Table("likes").Where("video_id = ?", videoId).Where("cancel = ? ", 1).Find(&likes)
	return likes, result.Error
}

func FindLikesByUserId(userId int64) ([]models.Like, error) {
	var likes []models.Like
	result := Db.Table("likes").Where("user_id = ?", userId).Where("cancel = ? ", 1).Find(&likes)
	return likes, result.Error
}
