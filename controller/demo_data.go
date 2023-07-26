package controller

import "github.com/RaymondCode/simple-demo/models"

var DemoVideos = []models.Video{
	{
		CommonEntity: models.NewCommonEntity(),
		//Id:            1,
		//Author:        DemoUser,
		PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

var DemoComments = []models.Comment{
	{
		CommonEntity: models.NewCommonEntity(),
		//Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = models.User{
	CommonEntity: models.NewCommonEntity(),
	//Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
}
