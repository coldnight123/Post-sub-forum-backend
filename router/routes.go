package router

import (
	"begin/goweb/controller"
	"begin/goweb/logger"
	"begin/goweb/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(),
		logger.GinRecovery(true),
		middlewares.RateLimitMiddleware(2*time.Second, 1),
	)
	v1 := r.Group("/api/v1")
	v1.POST("/signup", controller.SignupHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.CreatePostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		v1.GET("/posts2", controller.GetPostListHandler2)

		v1.POST("/vote", controller.PostVoteController)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
