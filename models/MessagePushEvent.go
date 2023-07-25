package models

type MessagePushEvent struct {
	CommonEntity
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
