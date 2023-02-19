package models

import "gorm.io/gorm"

type User struct {
	Id            int64  `json:"id,omitempty" gorm:"primarykey" `
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	Password      string `json:"password,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
}

type UserTable struct {
	gorm.Model
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	Password      string `json:"password,omitempty"`
}
