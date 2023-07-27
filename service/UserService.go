package service

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetUserById(Id int64) (models.User, error)

	GetUserByName(name string) (models.User, error)

	Save(user models.User) error

	Register(username string, password string, c *gin.Context) error

	Login(username string, password string, c *gin.Context) error

	UserInfo(token string) (*models.User, error)
}
