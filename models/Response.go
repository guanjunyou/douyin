package models

import "github.com/RaymondCode/simple-demo/utils"

type Response struct {
	utils.CommonEntity
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
