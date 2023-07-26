package controller

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	models.Response
	VideoList []models.VideoDVO `json:"video_list,omitempty"`
	NextTime  int64             `json:"next_time,omitempty"`
}

func GetVideoService() service.VideoServiceImpl {
	var videoService service.VideoServiceImpl
	return videoService
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	videoDVOList, err := GetVideoService().GetVideoList()
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  models.Response{StatusCode: 0},
		VideoList: videoDVOList,
		NextTime:  time.Now().Unix(),
	})
}
