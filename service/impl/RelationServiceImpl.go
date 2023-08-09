package impl

import (
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/mq"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/sirupsen/logrus"
	"sync"
)

type RelationServiceImpl struct {
	Logger *logrus.Logger
}

// FollowUser 关注用户
func (relationServiceImpl RelationServiceImpl) FollowUser(userId int64, toUserId int64, actionType int) error {
	relationServiceImpl.Logger.Info("FollowUser\n")
	if userId == toUserId {
		return fmt.Errorf("你不能关注(或者取消关注)自己")
	}
	followData := models.FollowMQToUser{
		UserId:       userId,
		FollowUserId: toUserId,
		ActionType:   actionType,
	}
	message, err := json.Marshal(followData)
	if err != nil {
		return err
	}
	mq.FollowRMQ.Publish(message)
	return nil
}

// GetFollows 查询关注列表
func (relationServiceImpl RelationServiceImpl) GetFollows(userId int64) ([]models.User, error) {
	relationServiceImpl.Logger.Info("GetFollows\n")
	var users []models.User
	err := utils.GetMysqlDB().Table("follow").Where("user_id = ? AND is_deleted != ?", userId, 1).Find(&users).Error
	if err != nil {
		return nil, err
	}

	//协程并发更新，isFollow 为 True 前端才能显示已关注
	var wg sync.WaitGroup
	for i := 0; i < len(users); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			users[i].IsFollow = true
		}(i)
	}

	wg.Wait()

	return users, nil
}

// GetFollowers 查询粉丝列表
func (relationServiceImpl RelationServiceImpl) GetFollowers(userId int64) ([]models.User, error) {
	relationServiceImpl.Logger.Info("GetFollowers")
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

func MakeFollowGroutine(count int) {
	for i := 0; i < count; i++ {
		go mq.FollowRMQ.Consumer()
	}
}
