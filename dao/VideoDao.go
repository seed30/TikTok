package dao

import "github.com/seed30/TikTok/models"

func GetAllVideos(videos *[]models.Video) {
	Db.Order("created_at desc").Preload("Author").Find(&videos)
}

func UpLoadInfoInsert(videos *models.VideoTable) error {
	//create := Db.Omit("IsFavorite").Create(&videos)
	create := Db.Table("videos").Create(&videos)
	return create.Error
}

func GetAllVideosByUserId(videos *[]models.Video, userId int64) {
	Db.Order("created_at desc").Preload("Author").Find(&videos, "author_id = ?", userId)
}

func UpdateFavoriteByVideoId(videoId int64, favoriteCount int64) {

	Db.Table("videos").Model(&models.VideoTable{}).Where("id = ?", videoId).Update("favorite_count", favoriteCount)
}

func UpdateCommentCountByVideoId(videoId int64, commentCount int64) {

	Db.Table("videos").Model(&models.VideoTable{}).Where("id = ?", videoId).Update("comment_count", commentCount)
}

func FindVideoByVideoId(videoId int64) models.VideoTable {
	var video models.VideoTable
	Db.Table("videos").First(&video, "id = ?", videoId)
	return video
}

func GetVideosByVideoIds(videoIds []int64, videos *[]models.Video) error {

	err := Db.Table("videos").Find(&videos, videoIds).Error
	return err
}
