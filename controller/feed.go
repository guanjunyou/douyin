package controller

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service/impl"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	models.Response
	VideoList []models.VideoDVO `json:"video_list,omitempty"`
	NextTime  int64             `json:"next_time,omitempty"`
}

// 拼装 VideoService
func GetVideoService() impl.VideoServiceImpl {
	var videoService impl.VideoServiceImpl
	var userService impl.UserServiceImpl
	var favoriteServer impl.FavoriteServiceImpl
	videoService.UserService = userService
	videoService.FavoriteService = favoriteServer
	return videoService
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	latestTimeStr := c.Query("latest_time")
	token := c.Query("token")
	var userId int64 = -1

	log.Printf("时间戳", latestTimeStr)
	var latestTime time.Time
	if latestTimeStr != "0" {
		me, _ := strconv.ParseInt(latestTimeStr, 10, 64)
		latestTime = time.Unix(me/1000, 0)
	} else {
		latestTime = time.Now()
	}
	log.Printf("获取到的时间 %v", latestTime)

	if token != "" {
		userClaims, err0 := utils.AnalyseToken(token)
		if err0 != nil {
			log.Println("解析token失败")
		}
		userId = userClaims.CommonEntity.Id
	}

	videoDVOList, nextTime, err := GetVideoService().GetVideoListByLastTime(latestTime, userId)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  models.Response{StatusCode: 0},
		VideoList: videoDVOList,
		NextTime:  nextTime.Unix(),
	})
}
