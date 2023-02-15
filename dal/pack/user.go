package pack

import (
	"context"
	"errors"
	"github.com/seed30/TikTok/dal/db"
	"github.com/seed30/TikTok/kitex_gen/user"
	"gorm.io/gorm"
)

func User(ctx context.Context, u *db.User, fromID int64) (*user.User, error) {
	if u == nil {
		return &user.User{}, nil
	}

	followCount := int64(u.FollowingCount)
	followerCount := int64(u.FollowerCount)

	// true->fromID已关注u.ID，false-fromID未关注u.ID
	isFollow := false
	relation, err := db.GetRelation(ctx, fromID, int64(u.ID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if relation != nil {
		isFollow = true
	}
	return &user.User{
		UserId:        int64(u.ID),
		Username:      u.UserName,
		FollowCount:   &followCount,
		FollowerCount: &followerCount,
		IsFollow:      isFollow,
	}, nil
}

// Users pack list of user info
func Users(ctx context.Context, us []*db.User, fromID int64) ([]*user.User, error) {
	users := make([]*user.User, 0)
	for _, u := range us {
		user2, err := User(ctx, u, fromID)
		if err != nil {
			return nil, err
		}

		if user2 != nil {
			users = append(users, user2)
		}
	}
	return users, nil
}
