package service

import "github.com/RaymondCode/simple-demo/models"

type UserService interface {
	GetUserById(Id int64) (models.User, error)

	Save(user models.User) error

	// Register 用户注册
	Register(username string, password string)
	// Login 用户登录
	Login(username string, password string)
}
