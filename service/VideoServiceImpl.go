package service

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/jinzhu/copier"
)

type VideoServiceImpl struct {
	UserService
}

func (videoService VideoServiceImpl) GetVideoList() ([]models.VideoDVO, error) {
	videolist, err := models.GetVideoList()
	VideoDVOList := make([]models.VideoDVO, len(videolist))
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
		var videoDVO models.VideoDVO
		err1 := copier.Copy(&videoDVO, &videolist[i])
		if err1 != nil {
			return nil, err1
		}
		videoDVO.Author = user
		VideoDVOList[i] = videoDVO
		//VideoDVOList = append(VideoDVOList, models.VideoDVO{
		//	CommonEntity:  videolist[i].CommonEntity,
		//	Author:        user,
		//	PlayUrl:       videolist[i].PlayUrl,
		//	CoverUrl:      videolist[i].CoverUrl,
		//	FavoriteCount: videolist[i].FavoriteCount,
		//	CommentCount:  videolist[i].CommentCount,
		//	IsFavorite:    videolist[i].IsFavorite,
		//})
	}
	return VideoDVOList, nil
}
