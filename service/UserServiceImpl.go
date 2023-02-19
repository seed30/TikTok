package service

import (
	"errors"

	"github.com/seed30/TikTok/dao"
	"github.com/seed30/TikTok/models"
)

type UserServiceImpl struct {
}

// 根据name 获取全部的User对象
func (usi *UserServiceImpl) GetTableUserByUsername(name string) (models.User, error) {
	user, err := dao.GetTabelUserByUsername(name)
	tableUsers := user

	return tableUsers, err
}

// 根据name 和加密的token 创建用户，并保存在数据库中
func (usi *UserServiceImpl) CreateTableUser(name string, password string) (models.User, error) {

	user := models.User{Name: name, FollowCount: 0, FollowerCount: 0, IsFollow: false, Password: password}
	userQuery, _ := dao.GetTabelUserByUsername(name)
	if user.Name == userQuery.Name {
		// 说明数据库中存在该用户
		return userQuery, errors.New("User already exist")
	} else {
		tableUser, err := dao.CreateTableUser(user)
		return tableUser, err
	}
}

// 根据token值查询数据库中是否存在
func (usi *UserServiceImpl) GetTableUserByToken(token string) (bool, models.User) {
	userByToken, user := dao.GetTabelUserByToken(token)
	byToken := userByToken
	if byToken > 0 {
		return true, user
	} else {
		return false, user
	}
}
