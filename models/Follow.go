package models

import "github.com/RaymondCode/simple-demo/utils"

// Follow 关注关系的item
type Follow struct {
	utils.CommonEntity
	UserId       int64 `json:"UserId"`
	FollowUserId int64 `json:"FollowUserId"`
}
