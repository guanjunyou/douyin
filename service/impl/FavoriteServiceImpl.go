package impl

import (
	"errors"
	"log"

	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

// "github.com/RaymondCode/simple-demo/service"

type FavoriteServiceImpl struct {
}

// LikeVedio 点赞或者取消点赞
func (fsi FavoriteServiceImpl) LikeVedio(userId int64, vedioId int64, actionType int) error {
	var err error
	tx := utils.GetMysqlDB().Begin()
	l := models.Like{
		UserId:  userId,
		VideoId: vedioId,
	}

	var faInDB *(models.Like)
	faInDB, err = l.FindByUserIdAndVedioId(tx)
	if err != nil {
		log.Printf("查询点赞记录发生异常 = %v", err)
		tx.Rollback()
		return err
	}

	if actionType == 1 {

		if faInDB.Id != 0 {
			log.Printf("该视频已点赞")
			tx.Rollback()
			err = errors.New("-1")
			return err
		}

		if err = l.Insert(tx); err != nil {
			log.Printf("添加点赞记录发生异常 = %v", err)
			tx.Rollback()
			return err
		}

		if err = findVedioAndUpdateFavoriteCount(tx, vedioId, 1); err != nil {
			log.Printf("修改视频点赞数量发生异常 = %v", err)
			tx.Rollback()
			return err
		}

	} else if actionType == 2 {

		if faInDB == nil || faInDB.Id == 0 {
			log.Printf("未找到要取消的点赞记录")
			err = errors.New("-2")
			tx.Rollback()
			return err
		}

		if err = faInDB.Delete(tx); err != nil {
			log.Printf("删除点赞记录发生异常 = %v", err)
			tx.Rollback()
			return err
		}

		if err = findVedioAndUpdateFavoriteCount(tx, vedioId, -1); err != nil {
			log.Printf("修改视频点赞数量发生异常 = %v", err)
			tx.Rollback()
			return err
		}

	}

	tx.Commit()
	return err
}

// findVedioAndUpdateFavoriteCount 修改视频的点赞数量，count 为 +-1
func findVedioAndUpdateFavoriteCount(tx *gorm.DB, vid int64, count int64) (err error) {
	var vInDB models.Video
	if err = tx.Model(&models.Video{}).Where("id = ? and is_deleted = 0", vid).Take(&vInDB).Error; err != nil {
		log.Printf("查询视频发生异常 = %v", err)
		return
	}
	if err = tx.Model(&models.Video{}).Where("id = ?", vid).Update("favorite_count", vInDB.FavoriteCount+count).Error; err != nil {
		log.Printf("修改视频点赞数量发生异常 = %v", err)
		return
	}
	return
}

func (fsi FavoriteServiceImpl) QueryVediosOfLike(userId int64) ([]models.LikeVedioListDVO, error) {
	var l models.Like
	var res []models.LikeVedioListDVO
	var err error
	res, err = l.GetLikeVedioListDVO(userId)

	if err != nil {
		return res, err
	}

	return res, err
}
