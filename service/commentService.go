package service

type commentService interface {

	// PostComments 登录用户对视频进行评论
	//actionType=1-发表评论 ，2-删除评论
	PostComments(userId int64, vedioId int64, actionType int, commentText string, commentId int64)

	// CommentList 查看视频的所有评论，按发布时间倒序
	CommentList(userId int64, vedioId int64)
}
