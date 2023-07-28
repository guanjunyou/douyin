package impl

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service"
)

type VideoServiceImpl struct {
	service.UserService
}

func (videoService VideoServiceImpl) GetVideoList() ([]models.VideoDVO, error) {
	videolist, err := models.GetVideoList()
	VideoDVOList := make([]models.VideoDVO, config.VideoCount)
	if err != nil {
		return nil, err
	}
	for i := range videolist {
		var userId = videolist[i].AuthorId

		//一定要通过videoService来调用 userSevice
		user, err := videoService.UserService.GetUserById(userId)
		if err != nil {
			return nil, err
		}
		VideoDVOList = append(VideoDVOList, models.VideoDVO{
			CommonEntity:  videolist[i].CommonEntity,
			Author:        user,
			PlayUrl:       videolist[i].PlayUrl,
			CoverUrl:      videolist[i].CoverUrl,
			FavoriteCount: videolist[i].FavoriteCount,
			CommentCount:  videolist[i].CommentCount,
			IsFavorite:    videolist[i].IsFavorite,
		})
	}
	return VideoDVOList, nil
}
