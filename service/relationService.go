package service

type relationService interface {
	//关注用户
	followUser(userId int64, toUserId int64, actionType int)

	//查询关注列表
	getFollows(userId int64)

	//查询粉丝列表
	getFollowers(userId int64)

	//查询好友列表
	getFriends(usrId int64)
}
