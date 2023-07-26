package service

type messageService interface {

	//发送消息
	sendMsg(userId int64, toUserId int64, actionType int, content string)

	//查看消息记录
	getHistoryOfChat(userId int64, toUserId int64)
}
