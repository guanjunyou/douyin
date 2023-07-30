package main

import (
	"github.com/RaymondCode/simple-demo/router"
	"github.com/RaymondCode/simple-demo/service/impl"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

var SF *utils.Snowflake

func main() {
	initDeps()
	go impl.RunMessageServer()

	r := gin.Default()
	r.Use(utils.RefreshHandler())
	r.Use(utils.AuthAdminCheck())
	// 创建一个 Snowflake 实例，并指定机器 ID
	SF = utils.NewSnowflake()
	router.InitRouter1(r)
	pprof.Register(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 加载项目依赖
func initDeps() {
	utils.InitFilter()
}
