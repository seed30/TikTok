package models

type Relation struct {
	Id       int64 `json:"id,omitempty" gorm:"primarykey" `
	UserId   int64 `json:"user_id,omitempty"`
	ToUserId int64 `json:"to_user_id,omitempty"` // 被关注方的ID
	Cancel   int   `json:"cancel,omitempty"`     // 1-关注 ， 2-取消关注
}
