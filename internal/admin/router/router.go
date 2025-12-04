package router

import (
	"merchant_api/internal/admin/controller"
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
			"service": "admin",
		})
	})

	// 初始化控制器
	authController := controller.NewAdminAuthController()

	// API 路由组
	api := r.Group("/mer_admin")
	{
		// 认证路由（无需登录）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)   // 登录
			auth.POST("/logout", authController.Logout) // 登出
		}

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AdminAuthMiddleware())
		{
			storeCategoryController := controller.NewStoreCategoryController()
			storeCategory := authorized.Group("/store_category")
			{
				storeCategory.POST("", storeCategoryController.Create)
				storeCategory.GET("", storeCategoryController.List)
				storeCategory.GET("/:id", storeCategoryController.Get)
				storeCategory.PUT("/:id", storeCategoryController.Update)
				storeCategory.DELETE("/:id", storeCategoryController.Delete)
			}
		}

	}

	return r
}
