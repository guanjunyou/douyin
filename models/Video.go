package models

import (
	"github.com/RaymondCode/simple-demo/config"
)

type Video struct {
	CommonEntity
	//Id            int64  `json:"id,omitempty"`
	AuthorId      int64  `json:"author_id"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type VideoDVO struct {
	CommonEntity
	Author        User
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
}

func (table *Video) TableName() string {
	return "video"
}

func GetVideoList() ([]VideoDVO, error) {
	videolist := make([]Video, config.VideoCount)
	VideoDVOList := make([]VideoDVO, config.VideoCount)
	err := DB.Where("is_deleted != ?", 1).Find(&videolist).Error
	if err != nil {
		return nil, err
	}
	for i := range videolist {
		userId := videolist[i].AuthorId
		user, _ := GetUserById(userId)
		VideoDVOList = append(VideoDVOList, VideoDVO{
			CommonEntity:  videolist[i].CommonEntity,
			Author:        user,
			PlayUrl:       videolist[i].PlayUrl,
			CoverUrl:      videolist[i].CoverUrl,
			FavoriteCount: videolist[i].FavoriteCount,
			CommentCount:  videolist[i].CommentCount,
			IsFavorite:    videolist[i].IsFavorite,
		})
	}

	return VideoDVOList, nil
}
