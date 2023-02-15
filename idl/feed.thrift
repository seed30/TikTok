namespace go feed

include "user.thrift"

struct Video {
    1: i64 id
    2: user.User author
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