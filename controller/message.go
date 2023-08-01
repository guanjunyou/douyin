package controller

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service/impl"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetMessageService() impl.MessageServiceImpl {
	return impl.MessageServiceImpl{}
}

type MessageListResponse struct {
	models.Response
	Data []models.MessageDVO `json:"message_list,omitempty"`
}

func errRespond(c *gin.Context, err error, statusCode int32, statusMsg string) bool {
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: statusCode, StatusMsg: statusMsg})
		return true
	}
	return false
}

func responseMessageList(c *gin.Context, messageList []models.MessageDVO) {
	c.JSON(http.StatusOK, MessageListResponse{Response: models.Response{StatusCode: 0, StatusMsg: "Message list success"}, Data: messageList})
}

// MessageAction no practical effect, just errRespond if token is valid
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	userClaim, err := utils.AnalyseToken(token)
	if errRespond(c, err, 1, "Token is invalid") {
		return
	}

	user, err := GetUserService().GetUserByName(userClaim.Name)
	if errRespond(c, err, 1, "User doesn't exist") {
		return
	}

	toUserIdInt64, err := strconv.ParseInt(toUserId, 10, 64)
	if errRespond(c, err, 1, "to_user_id is invalid") {
		return
	}

	err = GetMessageService().SendMsg(user.Id, toUserIdInt64, 1, content)
	if errRespond(c, err, 1, "Message send failed") {
		return
	}

	c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "Message send success"})
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

	userClaim, err := utils.AnalyseToken(token)
	if errRespond(c, err, 1, "Token is invalid") {
		return
	}

	user, err := GetUserService().GetUserByName(userClaim.Name)
	if errRespond(c, err, 1, "User doesn't exist") {
		return
	}

	toUserIdInt64, err := strconv.ParseInt(toUserId, 10, 64)
	if errRespond(c, err, 1, "to_user_id is invalid") {
		return
	}

	messageList, err := GetMessageService().GetHistoryOfChat(user.Id, toUserIdInt64)
	if errRespond(c, err, 1, "Message get failed") {
		return
	}
	responseMessageList(c, messageList)
}
