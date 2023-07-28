package models

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/utils"
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

// 在 model 层禁止操作除了数据库实体类外的其它类！ 禁止调用其它model或者service!
func GetVideoList() ([]Video, error) {
	videolist := make([]Video, config.VideoCount)
	err := utils.DB.Where("is_deleted != ?", 1).Find(&videolist).Error
	if err != nil {
		return nil, err
	}
	return videolist, nil
}

func SaveVedio(video *Video) error {
	err := utils.DB.Create(video).Error
	return err
}
