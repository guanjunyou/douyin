package impl

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/jinzhu/copier"
	"time"
)

type VideoServiceImpl struct {
	service.UserService
}

func (videoService VideoServiceImpl) GetVideoListByLastTime(latestTime time.Time) ([]models.VideoDVO, time.Time, error) {
	videolist, err := models.GetVideoListByLastTime(latestTime)
	size := len(videolist)
	VideoDVOList := make([]models.VideoDVO, size)
	if err != nil {
		return nil, time.Time{}, err
	}
	for i := range videolist {
		var userId = videolist[i].AuthorId

		//一定要通过videoService来调用 userSevice
		user, err1 := videoService.UserService.GetUserById(userId)
		if err1 != nil {
			return nil, time.Time{}, err1
		}
		var videoDVO models.VideoDVO
		err2 := copier.Copy(&videoDVO, &videolist[i])
		if err2 != nil {
			return nil, time.Time{}, err2
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
	nextTime := time.Now()
	if len(videolist) != 0 {
		nextTime = videolist[len(videolist)-1].CreateDate
	}
	return VideoDVOList, nextTime, nil
}
