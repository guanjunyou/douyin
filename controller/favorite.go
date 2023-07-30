package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/service/impl"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/RaymondCode/simple-demo/utils/resultutil"
	"github.com/gin-gonic/gin"
)

// 接收点赞的结构体
type FavoriteActionReq struct {
	Token      string `form:"token"`
	VideoId    string `form:"video_id"`    // 视频id
	ActionType string `form:"action_type"` // 1-点赞，2-取消点赞
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {

	var faReq FavoriteActionReq
	if err := c.ShouldBind(&faReq); err != nil {
		log.Printf("点赞操作，绑定参数发生异常：%v \n", err)
		resultutil.GenFail(c, "参数错误")
		return
	}
	fmt.Printf("参数 %+v \n", faReq)

	videoId, err := strconv.ParseInt(faReq.VideoId, 10, 64)

	if err != nil {
		log.Printf("点赞操作，videoId字符串转换发生异常 = %v", err)
		resultutil.GenFail(c, "参数错误")
		return
	}

	// 从Token中获取Uid
	var userClaim *utils.UserClaims
	userClaim, err = utils.AnalyseToken(faReq.Token)

	if err != nil {
		log.Printf("解析token发生异常 = %v", err)
		return
	}
	userId := userClaim.CommonEntity.Id

	var actionType int
	actionType, err = strconv.Atoi(faReq.ActionType)

	if err != nil {
		log.Printf("点赞操作，actionType字符串转换发生异常 = %v", err)
		resultutil.GenFail(c, "参数错误")
		return
	}

	var fs impl.FavoriteServiceImpl
	if err = fs.LikeVedio(userId, videoId, actionType); err != nil {
		log.Printf("点赞发生异常 = %v", err)
		if err.Error() == "-1" {
			resultutil.GenFail(c, "该视频已点赞")
			return
		}

		if err.Error() == "-2" {
			resultutil.GenFail(c, "未找到要取消的点赞记录")
			return
		}

		resultutil.GenFail(c, "点赞发生错误")
		return
	}

	resultutil.GenSuccessWithOutMsg(c)
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {

	userIdStr := c.Query("user_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)

	if err != nil {
		log.Printf("获取喜欢列表，userId字符串转换发生异常 = %v", err)
		resultutil.GenFail(c, "参数错误")
		return
	}

	var fs impl.FavoriteServiceImpl
	res, err := fs.QueryVediosOfLike(userId)

	if err != nil {
		log.Printf("获取喜欢列表，获取发生异常 = %v", err)
		resultutil.GenFail(c, "获取失败")
		return
	}

	c.JSON(http.StatusOK, models.VideoListResponse2{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: res,
	})
}
