package main

import (
	"context"
	"github.com/seed30/TikTok/cmd/feed/service"
	"github.com/seed30/TikTok/dal/pack"
	feed "github.com/seed30/TikTok/kitex_gen/feed"
	"github.com/seed30/TikTok/pkg/errno"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// GetUserFeed implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) GetUserFeed(ctx context.Context, req *feed.FeedRequest) (resp *feed.FeedResponse, err error) {
	var uid int64 = 0
	if *req.Token != "" {
		claim, err := Jwt.ParseToken(*req.Token)
		if err != nil {
			resp = pack.BuildVideoResp(err)
			return resp, nil
		} else {
			uid = claim.Id
		}
	}

	vis, nextTime, err := service.NewGetUserFeedService(ctx).GetUserFeed(req, uid)
	if err != nil {
		resp = pack.BuildVideoResp(err)
		return resp, nil
	}

	resp = pack.BuildVideoResp(errno.Success)
	resp.VideoList = vis
	resp.NextTime = nextTime
	return resp, nil
	return
}

// GetVideoById implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) GetVideoById(ctx context.Context, req *feed.IdRequest) (resp *feed.Video, err error) {
	// TODO: Your code here...
	return
}
