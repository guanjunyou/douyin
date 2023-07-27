package service

type favoriteService interface {

	// LikeVedio 点赞视频
	//actionType=1-点赞，2-取消点赞
	LikeVedio(userId int64, vedioId int64, actionType int)
	// QueryVediosOfLike 查询用户的所有点赞视频
	QueryVediosOfLike(userId int64)
}
