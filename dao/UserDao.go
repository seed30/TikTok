package dao

import (
	"github.com/seed30/TikTok/models"
)

func GetTabelUserByUsername(name string) (models.User, error) {
	tableUser := models.User{}
	result := Db.First(&tableUser, "name = ?", name)

	return tableUser, result.Error
}

func CreateTableUser(user models.User) (models.User, error) {

	reuslt := Db.Create(&user)
	return user, reuslt.Error
}

func GetTabelUserByToken(token string) (int64, models.User) {
	// TODO 类型转换优化 https://www.zhihu.com/question/449267385
	tableUser := models.UserTable{}
	first := Db.Table("users").First(&tableUser, "password = ?", token)
	user := models.User{
		Id:            int64(tableUser.ID),
		Name:          tableUser.Name,
		FollowCount:   tableUser.FollowCount,
		FollowerCount: tableUser.FollowerCount,
		IsFollow:      tableUser.IsFollow,
		Password:      tableUser.Password,
	}
	return first.RowsAffected, user
}

func GetFollowUsers(follwIds []int64) []models.User {
	var users []models.User
	Db.Table("users").Where(follwIds).Find(&users)
	return users
}
