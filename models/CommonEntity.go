package models

import "time"

type CommonEntity struct {
	CreateTime time.Time
	IsDelete   bool
}
