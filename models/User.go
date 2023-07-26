package models

import "github.com/RaymondCode/simple-demo/config"

type User struct {
	CommonEntity
	//Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

func (table *User) TableName() string {
	return "user"
}

func GetVideoList() ([]Video, error) {
	videolist := make([]Video, config.VideoCount)
	result := DB.Where("is_delete != ", 1).Find(&videolist)
	if result.Error != nil {
		return videolist, result.Error
	}
	return videolist, nil
}
