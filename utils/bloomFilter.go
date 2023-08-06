package utils

import (
	"github.com/RaymondCode/simple-demo/models"
	"github.com/bits-and-blooms/bloom/v3"
	"strconv"
)

var BloomFilter *bloom.BloomFilter

func InitBloomFilter() {
	BloomFilter = bloom.NewWithEstimates(10000000, 0.01)

	//加入全部评论的id
	comments := models.GetAllCommentDBs()
	for _, comment := range comments {
		BloomFilter.Add([]byte(strconv.Itoa(int(comment.Id))))
	}
	//加入全部视频的id
	videos, _ := models.GetAllExistVideo()
	for _, video := range videos {
		BloomFilter.Add([]byte(strconv.Itoa(int(video.Id))))
	}
}
