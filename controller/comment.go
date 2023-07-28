package controller

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service/impl"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetCommentService() impl.CommentServiceImpl {
	return impl.CommentServiceImpl{}
}

type CommentListResponse struct {
	models.Response
	CommentList []models.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	models.Response
	Comment models.Comment `json:"comment,omitempty"`
}

func ParseVideoId(c *gin.Context) int64 {
	video_id, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, models.Response{StatusCode: 1, StatusMsg: "video_id is invalid"})
		return -1
	}
	return video_id
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			comment := models.Comment{
				CommonEntity: utils.NewCommonEntity(),
				//Id:         1,
				User:    user,
				Content: text,
			}
			video_id := ParseVideoId(c)
			err := GetCommentService().PostComments(comment, video_id)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, models.Response{StatusCode: 1, StatusMsg: "Comment failed"})
				return
			}

			c.JSON(http.StatusOK, CommentActionResponse{Response: models.Response{StatusCode: 0, StatusMsg: "Comment success"},
				Comment: comment})
			return
		}
		c.JSON(http.StatusOK, models.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    models.Response{StatusCode: 0, StatusMsg: "Comment list"},
		CommentList: GetCommentService().CommentList(ParseVideoId(c)),
	})
}
