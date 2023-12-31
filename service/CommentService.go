package service

import "github.com/RaymondCode/simple-demo/models"

type CommentService interface {

	// PostComments 登录用户对视频进行评论
	//actionType=1-发表评论 ，2-删除评论
	PostComments(comment models.Comment, video_id int64) error

	DeleteComments(commentId int64) error

	// CommentList 查看视频的所有评论，按发布时间倒序
	CommentList(videoId int64) []models.Comment
}
