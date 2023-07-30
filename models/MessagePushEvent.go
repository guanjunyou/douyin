package models

import "github.com/RaymondCode/simple-demo/utils"

type MessagePushEvent struct {
	utils.CommonEntity
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

func (messagePushEvent *MessagePushEvent) TableName() string {
	return "message_push_event"
}

func SaveMessagePushEvent(messagePushEvent *MessagePushEvent) error {
	return utils.GetMysqlDB().Create(messagePushEvent).Error
}
