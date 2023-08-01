package service

import "github.com/RaymondCode/simple-demo/models"

type MessageService interface {

	// SendMsg 发送消息
	SendMsg(userId int64, toUserId int64, actionType int, content string) error

	// GetHistoryOfChat 查看消息记录
	GetHistoryOfChat(userId int64, toUserId int64) ([]models.MessageDVO, error)
}
