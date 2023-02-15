package service

import (
	"context"
	"github.com/seed30/TikTok/dal/db"
	"github.com/seed30/TikTok/dal/pack"

	"github.com/seed30/TikTok/kitex_gen/user"
)

type MGetUserService struct {
	ctx context.Context
}

// NewMGetUserService new MGetUserService
func NewMGetUserService(ctx context.Context) *MGetUserService {
	return &MGetUserService{ctx: ctx}
}

// MGetUser multiple get list of user info
func (s *MGetUserService) MGetUser(req *user.MGetUserRequest, fromID int64) ([]*user.User, error) {
	modelUsers, err := db.MGetUsers(s.ctx, req.UserIds)
	if err != nil {
		return nil, err
	}
	user, err := pack.Users(s.ctx, modelUsers, fromID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
