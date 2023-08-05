package models

import (
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

// 点赞接口的参数
type Like struct {
	utils.CommonEntity
	VideoId int64 `json:"videoId" gorm:"column:video_id"` //点赞的视频
	UserId  int64 `json:"userId" gorm:"column:user_id"`   //点赞的用户
}

type VideoListResponse2 struct {
	Response
	VideoList []LikeVedioListDVO `json:"video_list"`
}

type LikeVedioListDVO struct {
	Video
	Author *User `json:"author" gorm:"foreignKey:AuthorId"`
}

type LikeMQToVideo struct {
	UserId     int64 `json:"user_id"`
	VideoId    int64 `json:"video_id"`
	ActionType int   `json:"action_type"`
}

type LikeMQToUser struct {
	UserId     int64 `json:"user_id"`
	VideoId    int64 `json:"video_id"`
	AuthorId   int64 `json:"author_id"`
	ActionType int   `json:"action_type"`
}

// 表名
func (table *Like) TableName() string {
	return "like"
}

// Update 更新
func (l *Like) Update(tx *gorm.DB) (err error) {
	err = tx.Where("id = ?", l.Id).Updates(l).Error
	return
}

// Insert 插入记录
func (l *Like) Insert(tx *gorm.DB) (err error) {
	l.CommonEntity = utils.NewCommonEntity()
	err = tx.Create(l).Error
	return
}

// Delete 删除
func (l *Like) Delete(tx *gorm.DB) (err error) {
	l.IsDeleted = 1
	return l.Update(tx)
}

// FindByUserIdAndVedioId 通过userId和VedioId查找
func (l *Like) FindByUserIdAndVedioId() (res *Like, err error) {
	res = &Like{}
	err = utils.GetMysqlDB().Model(Like{}).Where("video_id = ? and user_id = ? and is_deleted = 0", l.VideoId, l.UserId).Find(res).Error
	return
}

func (l *Like) CountByUserIdAndVedioId(tx *gorm.DB) (res *Like, err error) {
	res = &Like{}
	err = tx.Model(Like{}).Where("video_id = ? and user_id = ? and is_deleted = 0", l.VideoId, l.UserId).Find(res).Error
	return
}

// GetLikeVedioListDVO 查询喜欢的视频列表
func (l *Like) GetLikeVedioListDVO(userId int64) ([]LikeVedioListDVO, error) {
	tx := utils.GetMysqlDB()
	var err error
	res := make([]LikeVedioListDVO, 0)
	err = tx.Table("`like` l").Select("v.*").Joins(`LEFT JOIN video v ON l.video_id = v.id`).Where("l.user_id = ? and l.is_deleted = 0", userId).Preload("Author").Find(&res).Error

	return res, err
}

// GetLikeVedioListDVO 查询喜欢的视频Id
func (l *Like) GetLikeVedioIdList(userId int64) ([]int64, error) {
	tx := utils.GetMysqlDB()
	var err error
	res := make([]int64, 0)
	err = tx.Table(l.TableName()).Select("video_id").Where("user_id = ? and is_deleted = 0", userId).Find(&res).Error

	return res, err
}
