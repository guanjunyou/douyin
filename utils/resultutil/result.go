package resultutil

import (
	"net/http"

	"github.com/RaymondCode/simple-demo/models"
	"github.com/gin-gonic/gin"
)

func gen(c *gin.Context, code int32, msg string) {
	c.JSON(http.StatusOK, models.Response{StatusCode: code, StatusMsg: msg})
}

// GenSuccessWithMsg 返回带消息的成功
func GenSuccessWithMsg(c *gin.Context, msg string) {
	gen(c, 0, msg)
}

// GenSuccessWithOutMsg 返回不带消息的成功
func GenSuccessWithOutMsg(c *gin.Context) {
	gen(c, 0, "")
}

// GenFail 返回失败
func GenFail(c *gin.Context, msg string) {
	gen(c, 400, msg)
}
