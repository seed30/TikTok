package models

type Like struct {
	Id      int64 `json:"id,omitempty" gorm:"primarykey" `
	UserId  int64 `json:"user_id,omitempty"`
	VideoId int64 `json:"video_id,omitempty"`
	Cancel  int   `json:"cancel,omitempty"`
}
