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

func parseVideoId(c *gin.Context) int64 {
	video_id, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, models.Response{StatusCode: 1, StatusMsg: "video_id is invalid"})
		return -1
	}
	return video_id
}

func parseCommetId(c *gin.Context) int64 {
	comment_id, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, models.Response{StatusCode: 1, StatusMsg: "comment_id is invalid"})
		return -1
	}
	return comment_id
}

// CommentAction comment action, 1 for post, 2 for delete
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	userClaim, err := utils.AnalyseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "Token is invalid"})
		return
	}
	user, err := GetUserService().GetUserByName(userClaim.Name)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	if actionType == "1" {
		text := c.Query("comment_text")
		utils.InitFilter()

		textAfterFilter := utils.Filter.Replace(text, '*')

		comment := models.Comment{
			CommonEntity: utils.NewCommonEntity(),
			//Id:         1,
			User:    user,
			Content: textAfterFilter,
		}

		video_id := parseVideoId(c)
		err := GetCommentService().PostComments(comment, video_id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, models.Response{StatusCode: 1, StatusMsg: "Comment failed"})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{Response: models.Response{StatusCode: 0, StatusMsg: ""},
			Comment: comment})
		return
	} else if actionType == "2" {
		comment_id := parseCommetId(c)
		err := GetCommentService().DeleteComments(comment_id)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, models.Response{StatusCode: 1, StatusMsg: "Delete comment failed"})
			return
		}
		c.JSON(http.StatusOK, models.Response{StatusCode: 0})
		return
	}
	c.JSON(http.StatusOK, models.Response{StatusCode: 0})
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    models.Response{StatusCode: 0, StatusMsg: "Comment list"},
		CommentList: GetCommentService().CommentList(parseVideoId(c)),
	})
}
