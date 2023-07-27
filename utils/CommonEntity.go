package utils

import (
	"time"
)

type CommonEntity struct {
	Id         int64     `json:"id,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	IsDeleted  int64     `json:"is_deleted"`
}

func NewCommonEntity() CommonEntity {
	sf := NewSnowflake()
	return CommonEntity{
		Id:         sf.NextID(),
		CreateTime: time.Now(),
		IsDeleted:  0,
	}
}
