package impl

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path/filepath"
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

// Publish 投稿接口
// TODO 借助redis协助实现feed流
func Publish(data *multipart.FileHeader, userId int64, title string, c *gin.Context) error {
	//文件名
	filename := filepath.Base(data.Filename)
	//将文件名拼接用户id
	finalName := fmt.Sprintf("%d_%s", userId, filename)
	//保存文件的路径，暂时保存在本队public文件夹下
	saveFile := filepath.Join("./public/", finalName)
	//保存视频在数据库中
	video := models.Video{
		CommonEntity: utils.NewCommonEntity(),
		AuthorId:     userId,
		PlayUrl:      saveFile,
		CoverUrl:     "",
		Title:        title,
	}
	err := models.SaveVedio(&video)
	if err != nil {
		return err
	}
	//保存视频在本地中
	if err = c.SaveUploadedFile(data, saveFile); err != nil {
		return err
	}
	return nil
}
