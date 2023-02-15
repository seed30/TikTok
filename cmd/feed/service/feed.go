package service

import (
	"context"
	"github.com/seed30/TikTok/dal/db"
	"github.com/seed30/TikTok/dal/pack"
	"github.com/seed30/TikTok/kitex_gen/feed"
	"time"
)

const (
	LIMIT = 30 // 单次返回最大视频数
)

type GetUserFeedService struct {
	ctx context.Context
}

// NewGetUserFeedService new GetUserFeedService
func NewGetUserFeedService(ctx context.Context) *GetUserFeedService {
	return &GetUserFeedService{ctx: ctx}
}

// GetUserFeed get feed info.
func (s *GetUserFeedService) GetUserFeed(req *feed.FeedRequest, fromID int64) (vis []*feed.Video, nextTime int64, err error) {
	videos, err := db.MGetVideos(s.ctx, LIMIT, req.LatestTime)
	if err != nil {
		return vis, nextTime, err
	}

	if len(videos) == 0 {
		nextTime = time.Now().UnixMilli()
		return vis, nextTime, nil
	} else {
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	}

	if vis, err = pack.Videos(s.ctx, videos, &fromID); err != nil {
		nextTime = time.Now().UnixMilli()
		return vis, nextTime, err
	}

	return vis, nextTime, nil
}
