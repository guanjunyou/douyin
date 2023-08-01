package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/service/impl"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	models.Response
	UserList []models.User `json:"user_list"`
}

var relationService service.RelationService = impl.RelationServiceImpl{}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	fmt.Println("RelationAction: ", token, toUserId, actionType)

	userClaims, _ := utils.AnalyseToken(token)
	toUserIdInt, _ := strconv.ParseInt(toUserId, 10, 64)
	actionTypeInt, _ := strconv.Atoi(actionType)

	err := relationService.FollowUser(userClaims.CommonEntity.Id, toUserIdInt, actionTypeInt)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  "",
	})
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)
	followUser, err := relationService.GetFollows(userIdInt)
	if err != nil {
		log.Printf("GetFollows fail")
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: followUser,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)
	followUser, err := relationService.GetFollowers(userIdInt)
	if err != nil {
		log.Printf("GetFollows fail")
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: followUser,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)
	followUser, err := relationService.GetFriends(userIdInt)
	if err != nil {
		log.Printf("GetFollows fail")
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: followUser,
	})
}
