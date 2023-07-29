package impl

import (
	"errors"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
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

func (userService UserServiceImpl) GetUserByName(name string) (models.User, error) {
	result, err := models.GetUserByName(name)
	if err != nil {
		log.Printf("方法GetUserById() 失败 %v", err)
		return result, err
	}
	return result, nil
}

func (userService UserServiceImpl) Save(user models.User) error {
	return models.SaveUser(user)
}

/*
（
已完成
*/
func (userService UserServiceImpl) Register(username string, password string, c *gin.Context) error {
	var userIdSequence = int64(1)
	_, errName := userService.GetUserByName(username)
	if errName == nil {
		c.JSON(http.StatusBadRequest, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "用户名重复"},
		})
		return nil
	}
	//var userRequest UserRequest
	//if err := c.ShouldBindJSON(&userRequest); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//username := userRequest.Username
	//password := userRequest.Password
	//加密
	encrypt, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password = string(encrypt)

	atomic.AddInt64(&userIdSequence, 1)
	newUser := models.User{
		CommonEntity: utils.NewCommonEntity(),
		Name:         username,
		Password:     password,
	}

	err := userService.Save(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "Cant not save the User!"},
		})
	} else {
		token, err1 := utils.GenerateToken(username, password, newUser.CommonEntity)
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
	return nil
}

/*
*
已完成
*/
func (userService UserServiceImpl) Login(username string, password string, c *gin.Context) error {

	_, err := userService.GetUserByName(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "用户不存在，请注册!"},
		})
		return nil
	}

	user, err1 := userService.GetUserByName(username)
	if err1 != nil {
		return err1
	}

	pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if pwdErr != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "密码错误！"},
		})
		return pwdErr
	}

	token, err2 := utils.GenerateToken(username, password, user.CommonEntity)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "生成token失败"},
		})
		return err2
	}

	err3 := utils.SaveTokenToRedis(user.Name, token, time.Duration(config.TokenTTL*float64(time.Second)))
	if err3 != nil {
		log.Printf("Fail : Save token in redis !")
		// TODO 开发完成后整理这个返回体 返回信息不能这么填
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "无法保存token 请检查redis连接"},
		})
		return err3
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: models.Response{StatusCode: 0, StatusMsg: "登录成功！"},
		UserId:   user.Id,
		Token:    token,
	})
	return nil
}

func (userService UserServiceImpl) UserInfo(userId int64, token string) (*models.User, error) {
	//userClaims, err := utils.AnalyseToken(token)
	//if err != nil || userClaims == nil {
	//	return nil, errors.New("用户未登录")
	//}
	user, err1 := userService.GetUserById(userId)
	if err1 != nil {
		return nil, errors.New("用户不存在！")
	}
	return &user, nil
}

type UserResponse struct {
	models.Response
	User models.User `json:"user"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	models.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}
