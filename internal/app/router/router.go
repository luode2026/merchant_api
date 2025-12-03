package router

import (
	"merchant_api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 应用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "app",
		})
	})

	// API 路由组
	api := r.Group("/api/app")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			// TODO: 添加注册、登录等接口
			_ = auth
		}

		// 需要认证的路由
		// authorized := api.Group("")
		// authorized.Use(middleware.JWTAuth())
		// {
		// 	// TODO: 添加需要认证的接口
		// }
	}

	return r
}
