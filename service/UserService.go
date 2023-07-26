package service

import "github.com/RaymondCode/simple-demo/models"

type UserService interface {
	GetUserById(Id int64) (models.User, error)
	//用户注册
	register(username string, password string)
	//用户登录
	login(username string, password string)
}
