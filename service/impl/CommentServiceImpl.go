package impl

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/utils"
	"sort"
)

type CommentServiceImpl struct {
}

func (commentService CommentServiceImpl) PostComments(comment models.Comment, video_id int64) error {
	var video models.Video

	err := utils.GetMysqlDB().Where("id = ? AND is_deleted != ?", video_id, 1).First(&video).Error
	if err != nil {
		return err
	}

	commentDB := comment.ToCommentDB()
	commentDB.VideoId = video_id
	return models.SaveComment(&commentDB)
}

// CommentList 查看视频的所有评论，按发布时间倒序
func (commentService CommentServiceImpl) CommentList(vedioId int64) []models.Comment {
	Comments := models.GetCommentByVideoId(vedioId)
	sort.Sort(models.ByCreateDate(Comments))
	return Comments
}

func (commentService CommentServiceImpl) DeleteComments(commentId int64) error {
	return models.DeleteComment(commentId)
}
