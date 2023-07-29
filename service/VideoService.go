package service

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"time"
)

type VideoService interface {
	// Feed
	// 通过传入时间戳，当前用户的id，返回对应的视频切片数组，以及视频数组中最早的发布时间
	Feed(lastTime time.Time, userId int64) ([]models.Video, time.Time, error)

	// GetVideo
	// 传入视频id获得对应的视频对象
	GetVideo(videoId int64, userId int64) (models.Video, error)

	// Publish
	// 将传入的视频流保存在文件服务器中，并存储在mysql表中
	// 5.23 加入title
	Publish(data *multipart.FileHeader, userId int64, title string, c *gin.Context) error

	// PublishList
	// 通过userId来查询对应用户发布的视频，并返回对应的视频切片数组
	PublishList(userId int64) ([]models.VideoDVO, error)

	// GetVideoIdList
	// 通过一个作者id，返回该用户发布的视频id切片数组
	GetVideoIdList(userId int64) ([]int64, error)

	// GetVideoList
	GetVideoListByLastTime(latestTime time.Time) ([]models.VideoDVO, time.Time, error)
}
