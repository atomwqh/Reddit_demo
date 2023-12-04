package routes

import (
	_ "main/docs"
	"time"

	"github.com/gin-contrib/pprof"

	// 千万不要忘了导入把你上一步生成的docs
	"main/controllers"
	"main/logger"
	"main/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置发布模式，终端不输出日志信息
	}

	controllers.InitTrans("zh")
	r := gin.New()
	// 添加限流中间件（每两秒添加一个令牌）
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api/v1")
	//注册业务
	v1.POST("/signup", controllers.SignUpHandler)
	// 登入业务
	v1.POST("/login", controllers.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件

	{
		v1.GET("/community", controllers.CommunityHandler)           // 社区列表
		v1.GET("/community/:id", controllers.CommunityDetailHandler) // 社区详情

		v1.POST("/post", controllers.CreatePostHandler)     // 创建帖子
		v1.GET("/post/:id", controllers.PostDetailHandler)  // 帖子详情
		v1.GET("/posts/", controllers.GetPostListHandler)   // 帖子列表
		v1.GET("/posts2/", controllers.GetPostListHandler2) // 升级版的查询帖子
		//v1.GET("/post/:communityid", controllers.GetCommunityPostListHandler) // 根据社区id查询帖子 在上面的接口中合并

		v1.POST("/vote", controllers.PostVoteController) // 投票

	}
	pprof.Register(r) // 注册pprof相关路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
