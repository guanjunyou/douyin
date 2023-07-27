package models

import "github.com/RaymondCode/simple-demo/utils"

type Comment struct {
	utils.CommonEntity
	//Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}
