package models

import (
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
	"strconv"
)

type Comment struct {
	utils.CommonEntity
	//Id         int64  `json:"id,omitempty"`
	User    User   `json:"user"`
	Content string `json:"content,omitempty"`
}

type ByCreateDate []Comment

func (a ByCreateDate) Len() int {
	return len(a)
}

func (a ByCreateDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCreateDate) Less(i, j int) bool {
	return a[i].CreateDate.Compare(a[j].CreateDate) > 0
}

// CommentDB是数据库储存的Entity
type CommentDB struct {
	utils.CommonEntity
	//Id         int64  `json:"id,omitempty"`
	UserId  int64  `json:"user_id"`
	VideoId int64  `json:"video_id"`
	Content string `json:"content,omitempty"`
}

func (comment *Comment) ToCommentDB() CommentDB {
	return CommentDB{
		CommonEntity: comment.CommonEntity,
		UserId:       comment.User.Id,
		VideoId:      -1,
		Content:      comment.Content,
	}
}

func (commentDB *CommentDB) TableName() string {
	return "comment"
}

func SaveComment(commentDB *CommentDB) error {
	videoID := commentDB.VideoId
	//comment_count++
	tx := utils.GetMysqlDB().Begin()

	err := tx.Model(&Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Create(commentDB).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func DeleteComment(commentId int64) error {
	//set is_deleted = 1
	tx := utils.GetMysqlDB().Begin()

	var commentDB CommentDB
	err := tx.Where("id = ?", commentId).First(&commentDB).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&Video{}).Where("id = ?", commentDB.VideoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&CommentDB{}).Where("id = ?", commentId).Update("is_deleted", 1).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetCommentByVideoId(videoId int64) []Comment {
	var comments []Comment
	var commentDBs []CommentDB
	// 找到对应User
	GetUser := func(Id int64) (User, error) {
		var user User

		err := utils.GetMysqlDB().Where("id = ? AND is_deleted != ?", Id, 1).First(&user).Error
		if err != nil {
			return user, err
		}
		return user, nil
	}

	err := utils.GetMysqlDB().Debug().Where("video_id = ? AND is_deleted != ?", strconv.Itoa(int(videoId)), 1).Find(&commentDBs).Error
	if err != nil {
		return comments
	}
	//change comment_count
	err = utils.GetMysqlDB().Model(&Video{}).Where("id = ?", videoId).Update("comment_count", len(commentDBs)).Error
	if err != nil {
	}

	for _, commentDB := range commentDBs {
		user, err := GetUser(commentDB.UserId)
		if err != nil {
			user = User{Name: "未知用户"}
		}
		comments = append(comments, Comment{
			CommonEntity: commentDB.CommonEntity,
			User:         user,
			Content:      commentDB.Content,
		})
	}
	return comments
}
