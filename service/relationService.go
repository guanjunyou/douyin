package service

type relationService interface {
	// FollowUser 关注用户
	FollowUser(userId int64, toUserId int64, actionType int)

	// GetFollows 查询关注列表
	GetFollows(userId int64)

	// GetFollowers 查询粉丝列表
	GetFollowers(userId int64)

	// GetFriends 查询好友列表
	GetFriends(usrId int64)
}
