package main

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/mq"
	"github.com/RaymondCode/simple-demo/router"
	"github.com/RaymondCode/simple-demo/service/impl"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/RaymondCode/simple-demo/utils/bloomFilter"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

var SF *utils.Snowflake

func main() {
	initDeps()
	config.ReadConfig()
	logrus.SetLevel(logrus.DebugLevel)
	go impl.RunMessageServer()
	r := gin.Default()
	r.Use(utils.RefreshHandler())
	r.Use(utils.AuthAdminCheck())
	// 创建一个 Snowflake 实例，并指定机器 ID
	SF = utils.NewSnowflake()
	router.InitRouter1(r)
	pprof.Register(r)
	utils.CreateGORMDB()
	bloomFilter.InitBloomFilter()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 加载项目依赖
func initDeps() {
	utils.InitFilter()

	mq.InitRabbitMQ()

	mq.InitLikeRabbitMQ()
	mq.InitCommentRabbitMQ()
	mq.InitFollowRabbitMQ()

	mq.InitFollowRabbitMQ()
	//impl.MakeFollowGroutine()

	mq.MakeLikeChannel()
	impl.MakeLikeGroutine()

	mq.MakeCommentChannel()
	impl.MakeCommentGoroutine()

	mq.MakeFollowChannel()
	impl.MakeFollowGroutine()

	controller.GetUserService().MakeLikeConsumers()
	controller.GetUserService().MakeFollowConsumers()
}
