package impl

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"mime/multipart"
	"path/filepath"
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

// Publish 投稿接口
// TODO 借助redis协助实现feed流
func (videoService VideoServiceImpl) Publish(data *multipart.FileHeader, userId int64, title string, c *gin.Context) error {
	//从title中过滤敏感词汇
	replaceTitle := utils.Filter.Replace(title, '#')
	//文件名
	filename := filepath.Base(data.Filename)
	////将文件名拼接用户id
	//finalName := fmt.Sprintf("%d_%s", userId, filename)
	////保存文件的路径，暂时保存在本队public文件夹下
	//saveFile := filepath.Join("./public/", finalName)
	//保存视频在本地中
	// if err = c.SaveUploadedFile(data, saveFile); err != nil {
	if err := utils.UploadToServer(data); err != nil {
		return err
	}
	user, err1 := models.GetUserById(userId)
	if err1 != nil {
		return nil
	}
	//保存视频在数据库中
	video := models.Video{
		CommonEntity: utils.NewCommonEntity(),
		AuthorId:     userId,
		PlayUrl:      "http://" + config.Config.VideoServer.Addr2 + "/" + filename,
		CoverUrl:     "",
		Title:        replaceTitle,
	}
	err := models.SaveVideo(&video)
	if err != nil {
		return err
	}
	//用户发布作品数加1
	user.WorkCount = user.WorkCount + 1
	models.UpdateUser(user)
	return nil
}

// PublishList  发布列表
func (videoService VideoServiceImpl) PublishList(userId int64) ([]models.VideoDVO, error) {
	videoList, err := models.GetVediosByUserId(userId)
	if err != nil {
		return nil, err
	}
	size := len(videoList)
	VideoDVOList := make([]models.VideoDVO, size)
	for i := range videoList {
		var userId = videoList[i].AuthorId
		//一定要通过videoService来调用 userSevice
		user, err1 := models.GetUserById(userId)
		if err1 != nil {
			return nil, err1
		}
		var videoDVO models.VideoDVO
		err := copier.Copy(&videoDVO, &videoList[i])
		if err != nil {
			return nil, err
		}
		videoDVO.Author = user
		VideoDVOList[i] = videoDVO
	}
	return VideoDVOList, nil
}
