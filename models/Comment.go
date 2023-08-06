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

type CommentMQToVideo struct {
	utils.CommonEntity
	ActionType int    `json:"action_type"`
	UserId     User   `json:"user"`
	VideoId    int64  `json:"video_id"`
	Content    string `json:"content"`
	CommentID  int64  `json:"id"`
}

func (comment *CommentMQToVideo) ToCommentDB() CommentDB {
	return CommentDB{
		CommonEntity: comment.CommonEntity,
		UserId:       comment.UserId.Id,
		VideoId:      comment.VideoId,
		Content:      comment.Content,
	}
}

func (comment *CommentMQToVideo) ToComment() Comment {
	return Comment{
		CommonEntity: comment.CommonEntity,
		User:         comment.UserId,
		Content:      comment.Content,
	}
}

// CommentDB是数据库储存的Entity
type CommentDB struct {
	utils.CommonEntity
	//Id         int64  `json:"id,omitempty"`
	UserId  int64  `json:"user_id"`
	VideoId int64  `json:"video_id"`
	Content string `json:"content,omitempty"`
}

func (comment *CommentDB) ToComment() Comment {
	user, _ := GetUserById(comment.UserId)
	return Comment{
		CommonEntity: comment.CommonEntity,
		User:         user,
		Content:      comment.Content,
	}
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

func GetCommentDBById(commentId int64) (CommentDB, error) {
	var commentDB CommentDB
	err := utils.GetMysqlDB().Where("id = ? AND is_deleted != ?", commentId, 1).First(&commentDB).Error
	if err != nil {
		return commentDB, err
	}
	return commentDB, nil
}

func GetCommentByVideoId(videoId int64) []Comment {
	var comments []Comment
	var commentDBs []CommentDB
	// 找到对应User

	err := utils.GetMysqlDB().Debug().Where("video_id = ? AND is_deleted != ?", strconv.Itoa(int(videoId)), 1).Find(&commentDBs).Error
	if err != nil {
		return comments
	}
	//change comment_count
	err = utils.GetMysqlDB().Model(&Video{}).Where("id = ?", videoId).Update("comment_count", len(commentDBs)).Error
	if err != nil {
	}

	for _, commentDB := range commentDBs {

		comments = append(comments, commentDB.ToComment())
	}
	return comments
}

func GetAllCommentDBs() []CommentDB {
	var commentDBs []CommentDB
	err := utils.GetMysqlDB().Where("is_deleted != ?", 1).Find(&commentDBs).Error
	if err != nil {
		return commentDBs
	}
	return commentDBs
}
