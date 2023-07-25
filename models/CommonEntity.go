package models

import "time"

type CommonEntity struct {
	Id         int64     `json:"id,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
	IsDelete   bool      `json:"is_delete"`
}
