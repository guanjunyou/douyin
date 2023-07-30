package models

import "github.com/RaymondCode/simple-demo/utils"

type MessageSendEvent struct {
	utils.CommonEntity
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type ByCreateTime []MessageSendEvent

func (a ByCreateTime) Len() int {
	return len(a)
}

func (a ByCreateTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCreateTime) Less(i, j int) bool {
	return a[i].CreateDate.Before(a[j].CreateDate)
}

func (messageSendEvent *MessageSendEvent) TableName() string {
	return "message_send_event"
}

func SaveMessageSendEvent(messageSendEvent *MessageSendEvent) error {
	return utils.GetMysqlDB().Create(messageSendEvent).Error
}

func FindMessageSendEventByUserIdAndToUserId(userId int64, toUserId int64) ([]MessageSendEvent, error) {
	var messageSendEvents []MessageSendEvent
	err := utils.GetMysqlDB().Where("user_id = ? AND to_user_id = ?", userId, toUserId).Find(&messageSendEvents).Error
	if err != nil {
		return nil, err
	}
	return messageSendEvents, nil
}
