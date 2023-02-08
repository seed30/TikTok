namespace go comment

enum ErrCode {
    SuccessCode                = 0
    ServiceErrCode             = 10001
    ParamErrCode               = 10002
    UserAlreadyExistErrCode    = 10003
    AuthorizationFailedErrCode = 10004
}

struct BaseResp {
    1: i64 status_code
    2: string status_message
    3: i64 service_time
}

struct User {
    1: i64 id
    2: string name
    3: i64 follow_count
    4: i64 follower_count
    5: bool is_follow
}

struct Comment {
    1: i64 comment_id
    2: User user
    3: string content
    4: i64 create_date
}

struct CommentActionRequest {
    1: string token
    2: i64 video_id
    3: i32 action_type
    4: string comment_text
    5: i64 comment_id
}

struct CommentActionResponse {
    1: i32 status_code
    2: string status_msg
    3: Comment comment
}

struct CommentListRequest {
    1: string token
    2: i64 video_id
}

struct CommentListResponse {
    1: i32 status_code
    2: string status_msg
    3: Comment comment_list
}

service CommentService {
    CommentActionResponse CommentAction(1: CommentActionRequest req)
    CommentListResponse CommentList(1: CommentListRequest req)
}
