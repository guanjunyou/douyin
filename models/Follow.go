package models

import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

// Follow 关注关系的item
type Follow struct {
	utils.CommonEntity
	UserId       int64 `json:"UserId"`
	FollowUserId int64 `json:"FollowUserId"`
}

type FollowMQToUser struct {
	UserId       int64 `json:"user_id"`
	FollowUserId int64 `json:"follow_user_id"`
	ActionType   int   `json:"action_type"`
}

// 表名
func (table *Follow) TableName() string {
	return "follow"
}

// Update 更新
func (f *Follow) Update(tx *gorm.DB) (err error) {
	err = tx.Where("id = ?", f.Id).Updates(f).Error
	return
}

// Insert 插入记录
func (f *Follow) Insert(tx *gorm.DB) (err error) {
	f.CommonEntity = utils.NewCommonEntity()
	userId := f.UserId
	toUserId := f.FollowUserId
	if userId == toUserId {
		return fmt.Errorf("你不能关注(或者取消关注)自己")
	}
	sql := "CALL addFollowRelation(?, ?)"
	var result int
	tx.Raw(sql, userId, toUserId).Scan(&result)
	if result != 0 {
		return errors.New("插入失败")
	}
	return
}

// Delete 删除
func (f *Follow) Delete(tx *gorm.DB) (err error) {
	f.CommonEntity = utils.NewCommonEntity()
	userId := f.UserId
	toUserId := f.FollowUserId
	if userId == toUserId {
		return fmt.Errorf("你不能关注(或者取消关注)自己")
	}
	sql := "CALL delFollowRelation(?, ?)"
	var result int
	tx.Raw(sql, userId, toUserId).Scan(&result)
	if result != 0 {
		return errors.New("删除失败")
	}
	return
}
