package impl

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/utils"
	"log"
)

type RelationServiceImpl struct {
}

// FollowResult 定义存储过程的返回值
type followResult struct {
	Cnt int
}

// FollowUser 关注用户
func (relationServiceImpl RelationServiceImpl) FollowUser(userId int64, toUserId int64, actionType int) error {
	var sql string
	// 1 关注 2 取消
	switch actionType {
	case 1:
		sql = "CALL addFollowRelation(?, ?)"
	case 2:
		sql = "CALL delFollowRelation(?, ?)"
	default:
		log.Println("非法actionType")
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
	return nil, nil
}

// GetFollowers 查询粉丝列表
func (relationServiceImpl RelationServiceImpl) GetFollowers(userId int64) ([]models.User, error) {
	return nil, nil
}

// GetFriends 查询好友列表
func (relationServiceImpl RelationServiceImpl) GetFriends(usrId int64) ([]models.User, error) {
	return nil, nil
}
