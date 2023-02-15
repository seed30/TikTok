package db

import (
	"context"
	"github.com/seed30/TikTok/pkg/consts"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	Author        User   `gorm:"foreignKey:AuthorID"`
	AuthorID      int    `gorm:"index:idx_authorid;not null"`
	PlayUrl       string `gorm:"type:varchar(255);not null"`
	CoverUrl      string `gorm:"type:varchar(255)"`
	FavoriteCount int    `gorm:"default:0"`
	CommentCount  int    `gorm:"default:0"`
	Title         string `gorm:"type:varchar(50);not null"`
}

func (v *Video) TableName() string {
	return consts.VideoTableName
}

// MGetVideos multiple get list of videos info
func MGetVideos(ctx context.Context, limit int, latestTime *int64) ([]*Video, error) {
	videos := make([]*Video, 0)

	if latestTime == nil || *latestTime == 0 {
		curTime := time.Now().UnixMilli()
		latestTime = &curTime
	}
	conn := DB.WithContext(ctx)

	if err := conn.Limit(limit).Order("update_time desc").Find(&videos, "update_time < ?", time.UnixMilli(*latestTime)).Error; err != nil {
		return nil, err
	}
	return videos, nil
}
