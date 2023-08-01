package service

import "github.com/RaymondCode/simple-demo/models"

type MessageService interface {

	// SendMessage SendMsg 发送消息
	SendMessage(userId int64, toUserId int64, content string) error

	// GetHistoryOfChat 查看消息记录
	GetHistoryOfChat(userId int64, toUserId int64) ([]models.MessageDVO, error)
}
