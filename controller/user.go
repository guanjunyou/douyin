package controller

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]models.User{
	"zhangleidouyin": {
		CommonEntity: models.NewCommonEntity(),
		//Id:            1,
		FollowCount:   10,
		FollowerCount: 5,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	models.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	models.Response
	User models.User `json:"user"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 拼装 UserService
func GetUserService() service.UserServiceImpl {
	var userService service.UserServiceImpl
	return userService
}

func Register(c *gin.Context) {
	//username := c.Query("username")
	//password := c.Query("password")
	var userRequest UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := userRequest.Username
	password := userRequest.Password
	//加密
	encrypt, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password = string(encrypt)

	atomic.AddInt64(&userIdSequence, 1)
	newUser := models.User{
		CommonEntity: models.NewCommonEntity(),
		Name:         username,
		Password:     password,
	}

	err := GetUserService().Save(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "Cant not save the User!"},
		})
	} else {
		token, err1 := models.GenerateToken(username, password, newUser.CommonEntity)
		if err1 != nil {
			log.Printf("Can not get the token!")
		}
		err2 := utils.SaveTokenToRedis(newUser.Name, token, time.Duration(config.TokenTTL*float64(time.Second)))
		if err2 != nil {
			log.Printf("Fail : Save token in redis !")
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: models.Response{StatusCode: 0},
				UserId:   newUser.Id,
				Token:    token,
			})
		}
	}

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: models.Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
