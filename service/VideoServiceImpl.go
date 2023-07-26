package service

import (
	"github.com/RaymondCode/simple-demo/models"
	"log"
)

type VideoServiceImpl struct {
}

func (videoService VideoServiceImpl) GetVideoList() ([]models.Video, error) {
	result, err := models.GetVideoList()
	if err != nil {
		log.Printf("方法GetProblemList() 失败 #{err}")
		return result, err
	}
	return result, nil
}
