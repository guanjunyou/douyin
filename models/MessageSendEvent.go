package models

import "github.com/RaymondCode/simple-demo/utils"

type MessageSendEvent struct {
	utils.CommonEntity
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
