package impl

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/utils"
)

type RelationServiceImpl struct {
}

// FollowResult 定义存储过程的返回值
type followResult struct {
	Cnt int
}

// FollowUser 关注用户
func (relationServiceImpl RelationServiceImpl) FollowUser(userId int64, toUserId int64, actionType int) error {
	if userId == toUserId {
		return fmt.Errorf("你不能关注(或者取消关注)自己")
	}
	var sql string
	// 1 关注 2 取消
	switch actionType {
	case 1:
		sql = "CALL addFollowRelation(?, ?)"
	case 2:
		sql = "CALL delFollowRelation(?, ?)"
	default:
		return fmt.Errorf("非法actionType")
	}
	var result followResult
	utils.GetMysqlDB().Raw(sql, userId, toUserId).Scan(&result)
	if actionType != 2 && result.Cnt != 0 {
		return fmt.Errorf("已经")
	}
	return nil
}

// GetFollows 查询关注列表
func (relationServiceImpl RelationServiceImpl) GetFollows(userId int64) ([]models.User, error) {
	var users []models.User
	err := utils.GetMysqlDB().Table("follow").Where("user_id = ? AND is_deleted != ?", userId, 1).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetFollowers 查询粉丝列表
func (relationServiceImpl RelationServiceImpl) GetFollowers(userId int64) ([]models.User, error) {
	var users []models.User
	err := utils.GetMysqlDB().Table("follow").Where("follow_user_id = ? AND is_deleted != ?", userId, 1).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetFriends 查询好友列表
func (relationServiceImpl RelationServiceImpl) GetFriends(userId int64) ([]models.User, error) {
	follows, err := relationServiceImpl.GetFollows(userId)
	if err != nil {
		return nil, err
	}
	followers, err := relationServiceImpl.GetFollowers(userId)
	if err != nil {
		return nil, err
	}
	var friends []models.User
	for _, user := range followers {
		if containsID(follows, user.Id) {
			friends = append(friends, user)
		}
	}
	return friends, nil
}

// containsID 辅助函数，用于检查指定的 id 是否在数组中存在
func containsID(arr []models.User, id int64) bool {
	for _, u := range arr {
		if u.Id == id {
			return true
		}
	}
	return false
}
