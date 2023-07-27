package models

import "github.com/RaymondCode/simple-demo/utils"

type Message struct {
	utils.CommonEntity
	//Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}
