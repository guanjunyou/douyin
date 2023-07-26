package models

import (
	"github.com/RaymondCode/simple-demo/utils"
	"time"
)

type CommonEntity struct {
	Id         int64     `json:"id,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	IsDelete   int64     `json:"is_delete"`
}

func NewCommonEntity() CommonEntity {
	sf := utils.NewSnowflake()
	return CommonEntity{
		Id:         sf.NextID(),
		CreateTime: time.Now(),
		IsDelete:   0,
	}
}
