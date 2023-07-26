package service

type favoriteService interface {

	//点赞视频
	//actionType=1-点赞，2-取消点赞
	likeVedio(userId int64, vedioId int64, actionType int)
	//查询用户的所有点赞视频
	queryVediosOfLike(userId int64)
}
