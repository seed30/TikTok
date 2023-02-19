package models

import (
	"time"
)

type Comment struct {
	Id          int64     `json:"id,omitempty" gorm:"primarykey" `
	User        User      `gorm:"foreignKey:UserId"`
	UserId      int64     `json:"user_id,omitempty"`
	VideoId     int64     `json:"video_id,omitempty"`
	Cancel      int       `json:"cancel,omitempty"` // 确保当前值是删除 还是未删除状态
	Content     string    `json:"content,omitempty"`
	CreatedDate time.Time `json:"created_date,omitempty"`
}
