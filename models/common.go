package models

import "errors"

var ErrCacheMiss = errors.New("chche missed")
var ErrLikeAction = errors.New("the like action is wrong")
var ErrVideoUserConsist = errors.New("the user data and video data in likes DB are not consistent")
var ErrRedis = errors.New("Redis has some error")
