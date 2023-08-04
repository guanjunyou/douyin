package impl

import (
	"encoding/json"
	"errors"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/mq"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
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
		token, err1 := utils.GenerateToken(username, newUser.CommonEntity)
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

	token, err2 := utils.GenerateToken(username, user.CommonEntity)
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

// LikeConsume  消费"userLikeMQ"中的消息
func (userService UserServiceImpl) LikeConsume(l *mq.LikeMQ) {
	_, err := l.Channel.QueueDeclare(l.QueueUserName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	//2、接收消息
	messages, err1 := l.Channel.Consume(
		l.QueueUserName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}
	go userService.likeConsume(messages)
	//forever := make(chan bool)
	//log.Println(messages)

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")
}

// 具体消费逻辑
// TODO  实现User与Video的解耦
func (userService UserServiceImpl) likeConsume(message <-chan amqp.Delivery) {
	for d := range message {
		jsonData := string(d.Body)
		log.Printf("user收到的消息为 %s\n", jsonData)
		data := models.LikeMQToVideo{}
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			panic(err)
		}
		userId := data.UserId
		//获得当前用户
		user, err := models.GetUserById(userId)
		videoId := data.VideoId
		//检索点赞视频
		video, err1 := models.GetVideoById(videoId)
		if err1 != nil {
			panic(err1)
		}
		//查询视频作者
		author, err2 := models.GetUserById(video.AuthorId)
		if err2 != nil {
			panic(err2)
		}
		actionType := data.ActionType
		tx := utils.GetMysqlDB().Begin()
		if actionType == 1 {
			//喜欢数量+一
			user.FavoriteCount = user.FavoriteCount + 1
			err = models.UpdateUser(tx, user)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}
			//总点赞数+1
			author.TotalFavorited = author.TotalFavorited + 1
			err = models.UpdateUser(tx, author)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}

		} else {
			//喜欢数量-1
			user.FavoriteCount = user.FavoriteCount - 1
			err = models.UpdateUser(tx, user)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}
			//总点赞数-1
			author.TotalFavorited = author.TotalFavorited - 1
			err = models.UpdateUser(tx, author)
			if err != nil {
				log.Println("err:", err)
				tx.Rollback()
				panic(err)
			}
		}
		tx.Commit()
	}
}

// 创建消费者协程
func (userService UserServiceImpl) MakeLikeConsumers() {
	numConsumers := 20
	for i := 0; i < numConsumers; i++ {
		go userService.LikeConsume(mq.LikeRMQ)
	}
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
