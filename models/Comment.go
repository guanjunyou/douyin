package models

import "github.com/RaymondCode/simple-demo/utils"

type Comment struct {
	utils.CommonEntity
	//Id         int64  `json:"id,omitempty"`
	User    User   `json:"user"`
	Content string `json:"content,omitempty"`
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
	return utils.DB.Create(commentDB).Error
}

func DeleteComment(commentId int64) error {
	//set is_deleted = 1
	return utils.DB.Model(&CommentDB{}).Where("id = ?", commentId).Update("is_deleted", 1).Error
}

func GetCommentByVideoId(videoId int64) []Comment {
	var comments []Comment
	var commentDBs []CommentDB
	// 找到对应User
	GetUser := func(Id int64) (User, error) {
		var user User

		err := utils.DB.Where("id = ? AND is_deleted != ?", Id, 1).First(&user).Error
		if err != nil {
			return user, err
		}
		return user, nil
	}

	err := utils.DB.Where("video_id = ? AND is_deleted != ?", videoId, 1).Find(&commentDBs).Error
	if err != nil {
		return comments
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
