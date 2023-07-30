package models

import "github.com/RaymondCode/simple-demo/utils"

type Message struct {
	utils.CommonEntity
	//Id         int64  `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
}

type MessageDVO struct {
	Id         int64  `json:"id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	UserId     int64  `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
}

func (message *Message) TableName() string {
	return "message"
}

func SaveMessage(message *Message) error {
	return utils.GetMysqlDB().Create(message).Error
}
