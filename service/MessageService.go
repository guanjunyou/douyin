package service

type messageService interface {

	// SendMsg 发送消息
	SendMsg(userId int64, toUserId int64, actionType int, content string)

	// GetHistoryOfChat 查看消息记录
	GetHistoryOfChat(userId int64, toUserId int64)
}
