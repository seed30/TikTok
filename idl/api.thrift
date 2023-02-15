namespace go api

struct BaseResp {
    1: i64 status_code
    2: string status_message
    3: i64 service_time
}

struct User {
    1: i64 user_id
    2: string username
    3: string avatar
    4: optional i64 follow_count
    5: optional i64 follower_count
    6: bool is_follow
}

struct CreateUserRequest {
    1: string username (api.form="username", api.vd="len($) > 0")
    2: string password (api.form="password", api.vd="len($) > 0")
}

struct CreateUserResponse {
    1: BaseResp base_resp
}

struct CheckUserRequest {
    1: string username (api.form="username", api.vd="len($) > 0")
    2: string password (api.form="password", api.vd="len($) > 0")
}

struct CheckUserResponse {
    1: BaseResp base_resp
}

struct Video {
    1: i64 id
    2: User author
    3: string play_url
    4: string cover_url
    5: i64 favorite_count
    6: i64 comment_count
    7: bool is_favorite
    8: string title
}

struct FeedRequest {
    1: optional i64 latest_time
    2: optional string token
}

struct FeedResponse {
    1: i32 status_code // 状态码，0-成功，其他值-失败
    2: optional string status_msg // 返回状态描述
    3: list<Video> video_list // 视频列表
    4: i64 next_time // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}

struct IdRequest {
    1: i64 video_id
    2: i64 search_id
}

service FeedService {
    FeedResponse GetUserFeed(1: FeedRequest req)
    Video GetVideoById(2: IdRequest req)
}

service ApiService {
    CreateUserResponse CreateUser(1: CreateUserRequest req) (api.post="/v2/user/register")
    CheckUserResponse CheckUser(1: CheckUserRequest req) (api.post="/v2/user/login")
    FeedResponse GetUserFeed(1: FeedRequest req) (api.get="/v2/feed")
}