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
	"sync"
	"time"
)

type VideoServiceImpl struct {
	service.UserService
	service.FavoriteService
}

func (videoService VideoServiceImpl) GetVideoListByLastTime(latestTime time.Time, userId int64) ([]models.VideoDVO, time.Time, error) {
	videolist, err := models.GetVideoListByLastTime(latestTime)
	size := len(videolist)
	var wg sync.WaitGroup
	VideoDVOList := make([]models.VideoDVO, size)
	if err != nil {
		return nil, time.Time{}, err
	}
	//for i := range videolist {
	//	var authorId = videolist[i].AuthorId
	//
	//	//一定要通过videoService来调用 userSevice
	//	user, err1 := videoService.UserService.GetUserById(authorId)
	//	if err1 != nil {
	//		return nil, time.Time{}, err1
	//	}
	//	var videoDVO models.VideoDVO
	//	err2 := copier.Copy(&videoDVO, &videolist[i])
	//	if err2 != nil {
	//		return nil, time.Time{}, err2
	//	}
	//	videoDVO.Author = user
	//	videoDVO.IsFavorite = videoService.FavoriteService.FindIsFavouriteByUserIdAndVideoId(userId, videoDVO.Id)
	//	VideoDVOList[i] = videoDVO
	//}
	var err0 error
	for i := range videolist {
		var authorId = videolist[i].AuthorId
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// 通过 videoService 来调用 userService
			user, err1 := videoService.UserService.GetUserById(authorId)
			if err1 != nil {
				err0 = err1
				return
			}
			var videoDVO models.VideoDVO
			err2 := copier.Copy(&videoDVO, &videolist[i])
			if err2 != nil {
				err0 = err2
				return
			}
			videoDVO.Author = user
			if userId != -1 {
				videoDVO.IsFavorite = videoService.FavoriteService.FindIsFavouriteByUserIdAndVideoId(userId, videoDVO.Id)
			} else {
				videoDVO.IsFavorite = false
			}
			VideoDVOList[i] = videoDVO
		}(i)
	}

	wg.Wait()
	if err0 != nil {
		return nil, time.Time{}, err0
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
	coverName, err := utils.UploadToServer(data)
	if err != nil {
		return err
	}
	user, err1 := models.GetUserById(userId)
	if err1 != nil {
		return nil
	}
	//将扩展名修改为.png并返回新的string作为封面文件名
	//ext := filepath.Ext(filename)
	//name := filename[:len(filename)-len(ext)]
	//coverName := name + ".png"
	//保存视频在数据库中
	video := models.Video{
		CommonEntity: utils.NewCommonEntity(),
		AuthorId:     userId,
		PlayUrl:      "http://" + config.Config.VideoServer.Addr2 + "/videos/" + filename,
		CoverUrl:     "http://" + config.Config.VideoServer.Addr2 + "/photos/" + coverName,
		Title:        replaceTitle,
	}
	err2 := models.SaveVideo(&video)
	if err2 != nil {
		return err2
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
	//创建多个协程并发更新
	var wg sync.WaitGroup
	//接收协程产生的错误
	var err0 error
	for i := range videoList {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var userId = videoList[i].AuthorId
			//一定要通过videoService来调用 userSevice
			user, err1 := models.GetUserById(userId)
			if err1 != nil {
				err0 = err1
			}
			var videoDVO models.VideoDVO
			err := copier.Copy(&videoDVO, &videoList[i])
			if err != nil {
				err0 = err1
			}
			videoDVO.Author = user
			VideoDVOList[i] = videoDVO
		}(i)
	}
	wg.Wait()
	//处理协程内的错误
	if err0 != nil {
		return nil, err0
	}
	return VideoDVOList, nil
}
