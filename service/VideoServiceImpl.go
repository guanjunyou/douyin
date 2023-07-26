package service

import (
	"github.com/RaymondCode/simple-demo/models"
	"log"
)

type VideoServiceImpl struct {
}

func (videoService VideoServiceImpl) GetVideoList() ([]models.VideoDVO, error) {
	result, err := models.GetVideoList()
	if err != nil {
		log.Printf("方法GetVideoList() 失败 %v", err)
		return result, err
	}
	return result, nil
}
