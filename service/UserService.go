package service

import "github.com/RaymondCode/simple-demo/models"

type UserService interface {
	GetUserById(Id int64) (models.User, error)

	GetUserByName(name string) (models.User, error)

	Save(user models.User) error
}
