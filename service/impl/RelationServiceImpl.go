package impl

type RelationServiceImpl struct {
}

// FollowUser 关注用户
func FollowUser(userId int64, toUserId int64, actionType int) error {
	return nil
}

// GetFollows 查询关注列表
func GetFollows(userId int64) error {
	return nil
}

// GetFollowers 查询粉丝列表
func GetFollowers(userId int64) error {
	return nil
}

// GetFriends 查询好友列表
func GetFriends(usrId int64) error {
	return nil
}
