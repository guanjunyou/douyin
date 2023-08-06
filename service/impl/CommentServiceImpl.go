package impl

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/mq"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

type CommentServiceImpl struct {
}

func (commentService CommentServiceImpl) PostComments(comment models.Comment, video_id int64) error {
	user, err := UserServiceImpl{}.GetUserById(comment.User.Id)
	if err != nil {
		return err
	}

	utils.BloomFilter.Add([]byte(strconv.Itoa(int(comment.Id))))

	toMQ := models.CommentMQToVideo{
		CommonEntity: comment.CommonEntity,
		ActionType:   1,
		UserId:       user,
		VideoId:      video_id,
		Content:      comment.Content,
		CommentID:    -1,
	}
	mq.CommentChannel <- toMQ
	return nil
}

// CommentList 查看视频的所有评论，按发布时间倒序
func (commentService CommentServiceImpl) CommentList(videoId int64) []models.Comment {
	rdb := utils.GetRedisDB()

	exist := utils.BloomFilter.Test([]byte(strconv.Itoa(int(videoId))))
	if !exist {
		return nil
	}

	//get by id
	commentID := rdb.ZRange(context.Background(), strconv.Itoa(int(videoId)), 0, -1)
	if len(commentID.Val()) == 0 {
		comments := models.GetCommentByVideoId(videoId)
		//save to redis
		for _, comment := range comments {
			commentJSON, err := json.Marshal(comment)
			if err != nil {
				log.Print(err)
				continue
			}
			rdb.Set(context.Background(), "comment:"+strconv.Itoa(int(comment.Id)), commentJSON, time.Hour*24)
			rdb.ZAdd(context.Background(), strconv.Itoa(int(videoId)), &redis.Z{
				Score:  float64(comment.CreateDate.Unix()),
				Member: strconv.Itoa(int(comment.Id)),
			})
		}
		return comments
	}

	var comments []models.Comment
	//从redis中读取评论实体
	for _, id := range commentID.Val() {
		commentJSON := rdb.Get(context.Background(), "comment:"+id)
		var comment models.Comment
		if commentJSON.Val() == "" {
			commentID, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				log.Print(err)
				continue
			}
			commentDB, err := models.GetCommentDBById(commentID)
			if err != nil {
				log.Print(err)
				continue
			}
			commentJSON, err := json.Marshal(commentDB.ToComment())
			rdb.Set(context.Background(), "comment:"+id, commentJSON, time.Hour*24)
			comments = append(comments, commentDB.ToComment())
			continue
		}

		err := json.Unmarshal([]byte(commentJSON.Val()), &comment)
		if err != nil {
			log.Print(err)
			continue
		}
		comments = append(comments, comment)
	}
	log.Print(comments)
	return comments
}

func (commentService CommentServiceImpl) DeleteComments(commentId int64) error {
	rdb := utils.GetRedisDB()

	exist := utils.BloomFilter.Test([]byte(strconv.Itoa(int(commentId))))
	if !exist {
		return errors.New("comment id not exist")
	}

	//check id exist
	commentExistKey := "commentID:" + strconv.Itoa(int(commentId))
	if (rdb.Exists(context.Background(), commentExistKey)).Val() == 0 {
		_, err := models.GetCommentDBById(commentId)
		if err != nil {
			return err
		}
	}
	//delete comment id
	rdb.Del(context.Background(), commentExistKey)

	toMQ := models.CommentMQToVideo{
		ActionType: 2,
		UserId:     models.User{},
		VideoId:    -1,
		Content:    "",
		CommentID:  commentId,
	}
	mq.CommentChannel <- toMQ
	return nil
}

func CommentActionConsumer() {
	for {
		select {
		case commentMQ := <-mq.CommentChannel:
			switch commentMQ.ActionType {
			case 1:
				//save comment
				commentDB := commentMQ.ToCommentDB()
				err := models.SaveComment(&commentDB)
				if err != nil {
					log.Print(err)
					continue
				}
				//save to redis
				rdb := utils.GetRedisDB()
				//set comment id
				commentExistKey := "commentID:" + strconv.Itoa(int(commentDB.Id))
				//set comment to video
				rdb.ZAdd(context.Background(), strconv.Itoa(int(commentDB.VideoId)), &redis.Z{
					Score:  float64(commentDB.CreateDate.Unix()),
					Member: commentDB.Id,
				})
				commentJSON, err := json.Marshal(commentMQ.ToComment())
				rdb.Set(context.Background(), commentExistKey, commentJSON, 0)
				if err != nil {
					log.Print(err)
					continue
				}
			case 2:
				commentDB, _ := models.GetCommentDBById(commentMQ.CommentID)
				videoID := commentDB.VideoId
				models.DeleteComment(commentMQ.CommentID)

				rdb := utils.GetRedisDB()
				rdb.Del(context.Background(), "comment:"+strconv.Itoa(int(commentMQ.CommentID)))
				rdb.ZRem(context.Background(), strconv.Itoa(int(videoID)), commentMQ.CommentID)

			default:
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
}

func MakeCommentGoroutine() {
	numConsumer := 20
	for i := 0; i < numConsumer; i++ {
		go CommentActionConsumer()
	}
}
