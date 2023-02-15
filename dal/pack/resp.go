package pack

import (
	"errors"
	"github.com/seed30/TikTok/kitex_gen/feed"
	"time"

	"github.com/seed30/TikTok/kitex_gen/user"
	"github.com/seed30/TikTok/pkg/errno"
)

// BuildBaseResp build baseResp from error
func BuildBaseResp(err error) *user.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *user.BaseResp {
	return &user.BaseResp{StatusCode: err.ErrCode, StatusMessage: err.ErrMsg, ServiceTime: time.Now().Unix()}
}

// BuildVideoResp build VideoResp from error
func BuildVideoResp(err error) *feed.FeedResponse {
	if err == nil {
		return videoResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return videoResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return videoResp(s)
}

func videoResp(err errno.ErrNo) *feed.FeedResponse {
	return &feed.FeedResponse{StatusCode: int32(err.ErrCode), StatusMsg: &err.ErrMsg}
}

//// BuildPublishResp build PublishResp from error
//func BuildPublishResp(err error) *publish.DouyinPublishActionResponse {
//	if err == nil {
//		return publishResp(errno.Success)
//	}
//
//	e := errno.ErrNo{}
//	if errors.As(err, &e) {
//		return publishResp(e)
//	}
//
//	s := errno.ErrUnknown.WithMessage(err.Error())
//	return publishResp(s)
//}
//
//func publishResp(err errno.ErrNo) *publish.DouyinPublishActionResponse {
//	return &publish.DouyinPublishActionResponse{StatusCode: int32(err.ErrCode), StatusMsg: &err.ErrMsg}
//}

//// BuildPublishResp build PublishResp from error
//func BuildPublishListResp(err error) *publish.DouyinPublishListResponse {
//	if err == nil {
//		return publishListResp(errno.Success)
//	}
//
//	e := errno.ErrNo{}
//	if errors.As(err, &e) {
//		return publishListResp(e)
//	}
//
//	s := errno.ErrUnknown.WithMessage(err.Error())
//	return publishListResp(s)
//}
//
//func publishListResp(err errno.ErrNo) *publish.DouyinPublishListResponse {
//	return &publish.DouyinPublishListResponse{StatusCode: int32(err.ErrCode), StatusMsg: &err.ErrMsg}
//}
