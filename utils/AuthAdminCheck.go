package utils

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 免登录接口列表
var notAuthArr = map[string]string{
	"/douyin/feed/":          "1",
	"/douyin/user/register/": "1",
	"/douyin/user/login/":    "1",
}

/*
*
token刷新
*/
func RefreshHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.获取token
		token := c.Query("token")
		//如果token为空则尝试从body中拿
		if token == "" {
			token = c.PostForm("token")
		}
		//2.判断是否携带token
		if token == "" {
			return
		}
		//3.解析token
		userClaims, err := AnalyseToken(token)
		if err != nil || userClaims == nil || userClaims.IsDeleted == 1 {
			return
		}
		//4.根据token查redis
		tokenFromRedis, err := GetTokenFromRedis(userClaims.Name)
		if tokenFromRedis == "" {
			//4.1 如果token可以被正确解析，重建redis缓存
			err := SaveTokenToRedis(userClaims.Name, token, time.Duration(config.TokenTTL*float64(time.Second)))
			if err != nil {

				c.JSON(http.StatusForbidden, gin.H{"StatusCode": "1", "StatusMsg": "服务器异常"})
				c.Abort()
				return
			}
			return
		}
		//6.刷新token的有效期
		err = RefreshToken(userClaims.Name, time.Duration(config.TokenTTL*float64(time.Second)))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"StatusCode": "1", "StatusMsg": "用户未登录"})
			return
		}
		c.Next()
	}
}

/*
*
登录校验
*/
func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.不用登录的接口直接放行
		//log.Println(c.Request.URL.Path)
		inWhite := notAuthArr[c.Request.URL.Path]
		if inWhite == "1" {
			return
		}
		//2.获取token
		token := c.Query("token")
		//如果token为空则尝试从body中拿
		if token == "" {
			token = c.PostForm("token")
		}
		userClaims, err := AnalyseToken(token)
		if err != nil || userClaims == nil || userClaims.IsDeleted == 1 {
			c.JSON(http.StatusOK, gin.H{"StatusCode": "1", "StatusMsg": "用户未登录"})
			//阻止该请求
			c.Abort()
			return
		}
		//3.根据token查redis
		tokenFromRedis, err := GetTokenFromRedis(userClaims.Name)
		if tokenFromRedis == "" || err != nil {
			c.JSON(http.StatusOK, gin.H{"StatusCode": "1", "StatusMsg": "用户未登录"})
			//阻止该请求
			c.Abort()
			return
		}
		c.Next()
	}
}

///*
//*
//鉴权
//*/
//func AuthAdminCheck(token string) error {
//	claims, err := AnalyseToken(token)
//	if err != nil || claims == nil {
//		log.Printf("Can not find this token !")
//		return errors.New("Can not find this token !")
//	}
//	return nil
//}
