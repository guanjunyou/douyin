package models

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	utils.CommonEntity
	//Id            int64  `json:"id,omitempty"`
	AuthorId      int64  `json:"author_id"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

type VideoDVO struct {
	utils.CommonEntity
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title,omitempty"`
}

func (table *Video) TableName() string {
	return "video"
}

// GetVideoListByLastTime 在 model 层禁止操作除了数据库实体类外的其它类！ 禁止调用其它model或者service!
func GetVideoListByLastTime(latestTime time.Time) ([]Video, error) {
	videolist := make([]Video, config.VideoCount)
	err := utils.GetMysqlDB().Where("is_deleted != ? AND create_date < ? ", 1, latestTime).Order("create_date desc").Limit(config.VideoCount).Find(&videolist).Error
	if err != nil {
		return nil, err
	}
	return videolist, nil
}

func SaveVideo(video *Video) error {
	err := utils.GetMysqlDB().Create(video).Error
	return err
}

// GetVediosByUserId 根据用户id查询发布的视频
func GetVediosByUserId(userId int64) ([]Video, error) {
	vedios := make([]Video, config.VideoCount)
	err := utils.GetMysqlDB().Where("author_id = ? AND is_deleted != ?", userId, 1).Find(&vedios).Error
	if err != nil {
		return nil, err
	}
	return vedios, nil
}

func GetVideoById(videoId int64) (Video, error) {
	var video Video
	err := utils.GetMysqlDB().First(&video, videoId).Error
	return video, err
}

func UpdateVideo(tx *gorm.DB, video Video) {
	tx.Save(&video)
}

func GetAllExistVideo() ([]Video, error) {
	var videos []Video
	err := utils.GetMysqlDB().Where("is_deleted != ?", 1).Find(&videos).Error
	return videos, err
}
