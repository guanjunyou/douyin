package service

type commentService interface {

	//登录用户对视频进行评论
	//actionType=1-发表评论 ，2-删除评论
	postComments(userId int64, vedioId int64, actionType int, commentText string, commentId int64)

	//查看视频的所有评论，按发布时间倒序
	commentList(userId int64, vedioId int64)
}
