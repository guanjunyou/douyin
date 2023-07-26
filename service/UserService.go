package service

import "github.com/RaymondCode/simple-demo/models"

type UserService interface {
	GetUserById(Id int64) (models.User, error)
}
