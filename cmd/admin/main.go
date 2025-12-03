package main

import (
	"fmt"
	"merchant_api/internal/admin/router"
	"merchant_api/pkg/config"
	"merchant_api/pkg/database"
	"merchant_api/pkg/logger"
	"merchant_api/pkg/redis"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 初始化日志
	if err := logger.InitLogger(
		cfg.Logger.Level,
		cfg.Logger.Format,
		cfg.Logger.Output,
		cfg.Logger.FilePath,
	); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}

	logger.Info("Admin 服务启动中...")

	// 初始化数据库
	if err := database.InitMySQL(cfg.Database.MySQL); err != nil {
		logger.Fatal(fmt.Sprintf("初始化数据库失败: %v", err))
	}
	logger.Info("数据库连接成功")

	// 初始化 Redis
	if err := redis.InitRedis(cfg.Redis); err != nil {
		logger.Fatal(fmt.Sprintf("初始化 Redis 失败: %v", err))
	}
	logger.Info("Redis 连接成功")

	// 设置 Gin 模式
	// gin.SetMode(cfg.Server.Admin.Mode)

	// 初始化路由
	r := router.SetupRouter()

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Admin.Port)
	logger.Info(fmt.Sprintf("Admin 服务启动在端口 %d", cfg.Server.Admin.Port))

	if err := r.Run(addr); err != nil {
		logger.Fatal(fmt.Sprintf("启动服务失败: %v", err))
	}
}
