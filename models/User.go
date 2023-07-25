package models

import "gorm.io/gorm"

type User struct {
	CommonEntity
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

func (table *User) TableName() string {
	return "problem_basic"
}

func GetProblemList() *gorm.DB {
	return nil
}
