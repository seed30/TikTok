package social

import (
	dbmodel "github.com/seed30/TikTok/dal/dao/model"
	"github.com/seed30/TikTok/kitex_gen/douyin/model"
	"github.com/seed30/TikTok/kitex_gen/douyin/social"
	. "github.com/seed30/TikTok/pkg/configs"
	"github.com/seed30/TikTok/pkg/constants"
	"github.com/seed30/TikTok/util"
	"context"
)

// MessageList implements the SocialServiceImpl interface.
func (s *Service) MessageList(ctx context.Context, req *social.DouyinMessageChatRequest) (resp *social.DouyinMessageChatResponse, err error) {
	resp = social.NewDouyinMessageChatResponse()
	userID, toUserID := req.GetBaseReq().GetUserId(), req.GetToUserId()
	messageList, err := s.dao.Message.QueryMessageList(ctx, userID, toUserID)
	if err != nil {
		Log.Errorf("query message list err: %v, fid: %v, tid: %v", err, userID, toUserID)
		return resp, constants.ErrQueryRecord
	}
	messageList2 := make([]*model.Message, len(messageList))
	for i := 0; i < len(messageList); i++ {
		message1, message2 := messageList[i], &model.Message{}
		message2.Id = message1.ID
		message2.FromUserId = message1.UID
		message2.ToUserId = message1.ToUID
		message2.Content = message1.Content
		timeStr := util.TimeToString(message1.CreatedAt)
		message2.CreateTime = &timeStr
		messageList2[i] = message2
	}
	resp.MessageList = messageList2
	return
}

// SendMessage implements the SocialServiceImpl interface.
func (s *Service) SendMessage(ctx context.Context, req *social.DouyinMessageActionRequest) (resp *social.DouyinMessageActionResponse, err error) {
	resp = social.NewDouyinMessageActionResponse()
	if req.GetActionType() != 1 {
		return resp, constants.ErrUnsupportedOperation
	}
	message := &dbmodel.Message{
		UID:     req.GetBaseReq().GetUserId(),
		ToUID:   req.GetToUserId(),
		Content: req.GetContent(),
	}
	err = s.dao.Message.CreateRecord(ctx, message)
	if err != nil {
		Log.Errorf("send message err: %v", err)
		return resp, constants.ErrCreateRecord
	}
	return
}
