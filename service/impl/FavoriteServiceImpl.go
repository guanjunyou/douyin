package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/mq"
	"log"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

// "github.com/RaymondCode/simple-demo/service"

type FavoriteServiceImpl struct {
}

// LikeVedio 点赞或者取消点赞
func (favoriteService FavoriteServiceImpl) LikeVideo(userId int64, videoId int64, actionType int) error {
	//分布式锁 不能让用户连续两次点赞或者取消同一个视频的请求进入
	userIdStr := strconv.FormatInt(userId, 10)
	videoIdStr := strconv.FormatInt(videoId, 10)

	lockKey := config.LikeLock + userIdStr + videoIdStr
	unLockKey := config.UnLikeLock + userIdStr + videoIdStr

	if actionType == 1 {
		isSuccess, _ := utils.GetRedisDB().SetNX(context.Background(), lockKey, "0", time.Duration(config.LikeLockTTL)*time.Second).Result()
		if isSuccess == false {
			log.Println("已点赞")
			return errors.New("-1")
		} else {
			utils.GetRedisDB().Del(context.Background(), unLockKey)
		}
	} else {
		isSuccess, _ := utils.GetRedisDB().SetNX(context.Background(), unLockKey, "0", time.Duration(config.LikeLockTTL)*time.Second).Result()
		if isSuccess == false {
			log.Println("已取消")
			return errors.New("-2")
		} else {
			utils.GetRedisDB().Del(context.Background(), lockKey)
		}
	}
	var err error
	//tx := utils.GetMysqlDB().Begin()
	l := models.Like{
		UserId:  userId,
		VideoId: videoId,
	}
	var faInDB *(models.Like)

	var isExists bool = false
	userLikeKey := config.LikeKey + userIdStr
	// 看看缓存中有没有这个集合
	exits, _ := utils.GetRedisDB().Exists(context.Background(), userLikeKey).Result()
	if exits != 0 {
		// 看看这个集合中有没有这个ID
		result, _ := utils.GetRedisDB().SIsMember(context.Background(), userLikeKey, videoIdStr).Result()
		isExists = result
		// 如果缓存里面的 Set 里面没有就要从数据库Z]里面查
		if !isExists {
			faInDB, err = l.FindByUserIdAndVedioId()
			if err != nil {
				log.Printf("查询点赞记录发生异常 = %v", err)
				return err
			}
			if faInDB.Id != 0 {
				isExists = true
			}
		}
	} else {
		// 缓存中没有则从数据库找
		faInDB, err = l.FindByUserIdAndVedioId()
		if err != nil {
			log.Printf("查询点赞记录发生异常 = %v", err)
			return err
		}

		if faInDB.Id != 0 {
			isExists = true
		}

	}

	if actionType == 1 {

		if isExists {
			log.Printf("该视频已点赞")
			//tx.Rollback()
			err = errors.New("-1")
			return err
		}

		//if err = findVideoAndUpdateFavoriteCount(tx, videoId, 1); err != nil {
		//	log.Printf("修改视频点赞数量发生异常 = %v", err)
		//	//tx.Rollback()
		//	return err
		//}
		// 查视频的作者
		author, errAuthorId := models.GetVideoById(videoId)
		if errAuthorId != nil {
			log.Println("不能找到这个作者")
			return errAuthorId
		}
		authorId := author.AuthorId

		//mqData := models.LikeMQToVideo{UserId: userId, VideoId: videoId, ActionType: actionType}
		mqData := models.LikeMQToUser{UserId: userId, VideoId: videoId, ActionType: actionType, AuthorId: authorId}
		// 加入 channel
		mq.LikeChannel <- mqData
		jsonData, err := json.Marshal(mqData)
		if err != nil {
			log.Println("json序列化失败 = #{err}")
			//TODO 处理失败导致的数据不一致
		}
		//加入消息队列
		mq.LikeRMQ.Publish(string(jsonData))

		return nil

	} else if actionType == 2 {

		if !isExists && (faInDB == nil || faInDB.Id == 0) {
			log.Printf("未找到要取消的点赞记录")
			err = errors.New("-2")
			//tx.Rollback()
			return err
		}

		//if err = faInDB.Delete(tx); err != nil {
		//	log.Printf("删除点赞记录发生异常 = %v", err)
		//	tx.Rollback()
		//	return err
		//}
		//
		//if err = findVideoAndUpdateFavoriteCount(tx, videoId, -1); err != nil {
		//	log.Printf("修改视频点赞数量发生异常 = %v", err)
		//	tx.Rollback()
		//	return err
		//}
		// 查视频的作者
		author, errAuthorId := models.GetVideoById(videoId)
		if errAuthorId != nil {
			log.Println("不能找到这个作者")
			return errAuthorId
		}
		authorId := author.AuthorId

		//mqData := models.LikeMQToVideo{UserId: userId, VideoId: videoId, ActionType: actionType}
		mqData := models.LikeMQToUser{UserId: userId, VideoId: videoId, ActionType: actionType, AuthorId: authorId}
		// 加入 channel
		mq.LikeChannel <- mqData
		jsonData, err := json.Marshal(mqData)
		if err != nil {
			log.Printf("json序列化失败 = #{err}")
			//TODO 处理失败导致的数据不一致
		}
		//加入消息队列
		mq.LikeRMQ.Publish(string(jsonData))
		// TODO 消息队列处理失败会导致数据不一致

		return nil

	}

	//tx.Commit()
	return err
}

// findVideoAndUpdateFavoriteCount 修改视频的点赞数量，count 为 +-1
func findVideoAndUpdateFavoriteCount(tx *gorm.DB, vid int64, count int64) (err error) {
	var vInDB models.Video
	if err = tx.Model(&models.Video{}).Where("id = ? and is_deleted = 0", vid).Take(&vInDB).Error; err != nil {
		log.Printf("查询视频发生异常 = %v", err)
		return
	}
	fmt.Println(vInDB.CreateDate)
	if err = tx.Model(&models.Video{}).Where("id = ?", vid).Update("favorite_count", vInDB.FavoriteCount+count).Error; err != nil {
		log.Printf("修改视频点赞数量发生异常 = %v", err)
		return
	}
	fmt.Println(vInDB.CreateDate)
	return
}

func (favoriteService FavoriteServiceImpl) QueryVideosOfLike(userId int64) ([]models.LikeVedioListDVO, error) {
	likeKey := config.LikeKey + strconv.FormatInt(userId, 10)
	exists, _ := utils.GetRedisDB().Exists(context.Background(), likeKey).Result()
	if exists != 0 {
		likeIdsSet, _ := utils.GetRedisDB().SMembers(context.Background(), likeKey).Result()
		var res []models.LikeVedioListDVO

		for i := range likeIdsSet {
			id, _ := strconv.ParseInt(likeIdsSet[i], 10, 64)
			video, _ := models.GetVideoById(id)
			author, _ := models.GetUserById(video.AuthorId)
			var likeVideoListDVO models.LikeVedioListDVO
			likeVideoListDVO.Author = &author
			likeVideoListDVO.Video = video
			res = append(res, likeVideoListDVO)
		}
		return res, nil
	}

	var l models.Like

	var res []models.LikeVedioListDVO
	var err error
	res, err = l.GetLikeVedioListDVO(userId)

	_ = BuildLikeRedis(userId)
	if err != nil {
		return res, err
	}

	return res, err
}

func (favoriteService FavoriteServiceImpl) FindIsFavouriteByUserIdAndVideoId(userId int64, videoId int64) bool {
	//tx := utils.GetMysqlDB()
	likeKey := config.LikeKey + strconv.FormatInt(userId, 10)
	videoIdStr := strconv.FormatInt(videoId, 10)
	exists, _ := utils.GetRedisDB().Exists(context.Background(), likeKey).Result()
	if exists != 0 {
		videoExists, _ := utils.GetRedisDB().SIsMember(context.Background(), likeKey, videoIdStr).Result()
		return videoExists
	}
	like := models.Like{
		UserId:  userId,
		VideoId: videoId,
	}

	isLike, _ := like.FindByUserIdAndVedioId()

	if isLike.Id != 0 {
		return true
	} else {
		return false
	}
}

func LikeConsumer(ch <-chan models.LikeMQToUser) {
	for {
		select {
		case msg := <-ch:
			// 在这里处理接收到的消息
			tx := utils.GetMysqlDB().Begin()
			if msg.ActionType == 1 {
				like := models.Like{
					CommonEntity: utils.NewCommonEntity(),
					UserId:       msg.UserId,
					VideoId:      msg.VideoId,
				}
				err1 := like.Insert(tx)
				if err1 != nil {
					log.Printf(err1.Error())
					tx.Rollback()
				}
				video, err2 := models.GetVideoById(msg.VideoId)
				if err2 != nil {
					log.Printf(err2.Error())
					tx.Rollback()
				}
				video.FavoriteCount++
				models.UpdateVideo(tx, video)
				tx.Commit()

				userIdStr := strconv.FormatInt(msg.UserId, 10)
				videoIdStr := strconv.FormatInt(msg.VideoId, 10)
				likeSetKey := config.LikeKey + userIdStr
				exists, _ := utils.GetRedisDB().Exists(context.Background(), likeSetKey).Result()
				if exists == 0 { //缓存里面没有
					errBuildRedis := BuildLikeRedis(msg.UserId)
					if errBuildRedis != nil {
						log.Println("重建缓存失败", errBuildRedis)
					}
				} else {
					utils.GetRedisDB().SAdd(context.Background(), likeSetKey, videoIdStr)
				}
			}

			if msg.ActionType == 2 {
				like := models.Like{
					CommonEntity: utils.NewCommonEntity(),
					UserId:       msg.UserId,
					VideoId:      msg.VideoId,
				}
				findLike, err := like.FindByUserIdAndVedioId()
				if err != nil {
					tx.Rollback()
				}
				err1 := findLike.Delete(tx)
				if err1 != nil {
					log.Printf(err1.Error())
					tx.Rollback()
				}
				video, err2 := models.GetVideoById(msg.VideoId)
				if err2 != nil {
					log.Printf(err2.Error())
					tx.Rollback()
				}
				//TODO 防止减到负数
				video.FavoriteCount--
				models.UpdateVideo(tx, video)
				tx.Commit()

				userIdStr := strconv.FormatInt(msg.UserId, 10)
				videoIdStr := strconv.FormatInt(msg.VideoId, 10)
				likeSetKey := config.LikeKey + userIdStr
				// 删除缓存 SET 中的ID ， 避免脏数据产生
				utils.GetRedisDB().SRem(context.Background(), likeSetKey, videoIdStr)
			}
		default:
			// 如果channel为空，暂停一段时间后重新监听
			time.Sleep(time.Millisecond * 1)
		}
	}
}

// 重建缓存
func BuildLikeRedis(userId int64) error {
	like := models.Like{}
	idSet, err := like.GetLikeVedioIdList(userId)
	if err != nil {
		return err
	}
	userIdStr := strconv.FormatInt(userId, 10)
	likeSetKey := config.LikeKey + userIdStr
	var strValues []string
	for i := range idSet {
		strValues = append(strValues, strconv.FormatInt(idSet[i], 10))
	}

	ctx := context.Background()
	err = utils.GetRedisDB().SAdd(ctx, likeSetKey, strValues).Err()
	errSetTime := utils.GetRedisDB().Expire(ctx, likeSetKey, time.Duration(config.LikeKeyTTL)*time.Second).Err()
	if errSetTime != nil {
		log.Println("redis 时间设置失败", errSetTime.Error())
	}
	return err
}

// 创建消费者协程
func MakeLikeGroutine() {
	numConsumers := 20
	for i := 0; i < numConsumers; i++ {
		go LikeConsumer(mq.LikeChannel)
	}
}
