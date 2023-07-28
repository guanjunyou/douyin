package utils

import (
	"time"
)

type CommonEntity struct {
	Id         int64     `json:"id,omitempty"`
	CreateDate time.Time `json:"create_date,omitempty"`
	IsDeleted  int64     `json:"is_deleted"`
}

func NewCommonEntity() CommonEntity {
	sf := NewSnowflake()
	return CommonEntity{
		Id:         sf.NextID(),
		CreateDate: time.Now(),
		IsDeleted:  0,
	}
}
