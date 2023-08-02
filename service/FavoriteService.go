package service

import "github.com/RaymondCode/simple-demo/models"

type FavoriteService interface {

	// LikeVideo  点赞视频
	//actionType=1-点赞，2-取消点赞
	LikeVideo(userId int64, vedioId int64, actionType int) error
	// QueryVideosOfLike  查询用户的所有点赞视频
	QueryVideosOfLike(userId int64) ([]models.LikeVedioListDVO, error)

	FindIsFavouriteByUserIdAndVideoId(userId int64, videoId int64) bool
}
