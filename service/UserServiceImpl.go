package service

import (
	"github.com/RaymondCode/simple-demo/models"
	"log"
)

type UserServiceImpl struct {
}

func (userService UserServiceImpl) GetUserById(Id int64) (models.User, error) {
	result, err := models.GetUserById(Id)
	if err != nil {
		log.Printf("方法GetUserById() 失败 %v", err)
		return result, err
	}
	return result, nil
}

func (userService UserServiceImpl) Save(user models.User) error {
	return models.SaveUser(user)
}
