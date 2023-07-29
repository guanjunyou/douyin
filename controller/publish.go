package controller

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service/impl"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type VideoListResponse struct {
	models.Response
	VideoList []models.VideoDVO `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	//1.获取token并解析出user_id、data、title
	token := c.PostForm("token")
	userClaims, _ := utils.AnalyseToken(token)
	userId := userClaims.CommonEntity.Id
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title := c.Query("title")
	//2. 调用service层处理业务逻辑
	err = impl.Publish(data, userId, title, c)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		StatusCode: 0,
		StatusMsg:  "投稿成功！",
	})
	return
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	//获取用户id
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: models.Response{
				StatusCode: 1,
				StatusMsg:  "类型转换错误",
			},
			VideoList: nil,
		})
	}
	publishList, err := impl.PublishList(userId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: models.Response{
				StatusCode: 1,
				StatusMsg:  "数据库异常",
			},
			VideoList: nil,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "查询成功",
		},
		VideoList: publishList,
	})
}
