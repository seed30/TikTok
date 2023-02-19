include "base.thrift"
include "social.thrift"

namespace go douyin.api

service DouyinApi {
    // basic service
    base.douyin_feed_response Feed(1:base.douyin_feed_request req)(api.get="/douyin/feed/")
    base.douyin_user_register_response UserRegister(1:base.douyin_user_register_request req)(api.post="/douyin/user/register/")
    base.douyin_user_login_response UserLogin(1:base.douyin_user_login_request req)(api.post="/douyin/user/login/")
    base.douyin_user_response UserMsg(1:base.douyin_user_request req)(api.get="/douyin/user/")
    base.douyin_publish_action_response PublishAction(1:base.douyin_publish_action_request req)(api.post="/douyin/publish/action/")
    base.douyin_publish_list_response PublishList(1:base.douyin_publish_list_request req)(api.get="/douyin/publish/list/")

    // social service
    social.douyin_follow_action_response FollowAction(1:social.douyin_follow_action_request req)(api.get="/douyin/relation/action/")
    social.douyin_following_list_response FollowList(1:social.douyin_following_list_request req)(api.get="/douyin/relatioin/follow/list/")
    social.douyin_follower_list_response FollowerList(1:social.douyin_follower_list_request req)(api.get="/douyin/relation/follower/list/")
    social.douyin_relation_friend_list_response FriendList(1:social.douyin_relation_friend_list_request req)(api.get="/douyin/relation/friend/list/")
    social.douyin_message_chat_response MessageList(1:social.douyin_message_chat_request req)(api.get="/douyin/message/chat/")
    social.douyin_message_action_response SendMessage(1:social.douyin_message_action_request req)(api.post="/douyin/message/action/")
}