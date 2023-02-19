package service

import (
	"github.com/seed30/TikTok/dao"
	"github.com/seed30/TikTok/models"
)

type VideoServiceImpl struct {
}

func (vsi *VideoServiceImpl) GetAllVideos() []models.Video {
	var videos []models.Video
	dao.GetAllVideos(&videos)

	return videos
}
func (vsi *VideoServiceImpl) GetAllVideosByUserId(userId int64) []models.Video {
	var videos []models.Video
	dao.GetAllVideosByUserId(&videos, userId)
	return videos
}

func (vsi *VideoServiceImpl) UpdateFavoriteByVideoId(videoId int64) {

	favoriteCount := dao.SumFavoriteCountByVideoId(videoId)

	dao.UpdateFavoriteByVideoId(videoId, favoriteCount)
}

func (vsi *VideoServiceImpl) UpdateCommentCountByVideoId(videoId int64) {
	commentCount := dao.SumCommentCountByVideoId(videoId)

	dao.UpdateCommentCountByVideoId(videoId, commentCount)
}

func (vsi *VideoServiceImpl) FindVideoByVideoId(videoId int64) models.VideoTable {

	videoTable := dao.FindVideoByVideoId(videoId)
	return videoTable

}

func (vsi *VideoServiceImpl) GetVideosByVideoIds(videoIds []int64) ([]models.Video, error) {
	var videos []models.Video
	err := dao.GetVideosByVideoIds(videoIds, &videos)

	return videos, err
}
