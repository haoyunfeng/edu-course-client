package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haoyunfeng/edu-course-client/internal/config"
	"github.com/haoyunfeng/edu-course-client/internal/router"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置 gin 模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 gin 路由
	r := gin.Default()

	// 统一设置所有路由和服务（只需一行代码）
	services, err := router.SetupRoutes(r, cfg)
	if err != nil {
		log.Fatalf("Failed to setup routes: %v", err)
	}
	defer services.Close()

	// 启动 HTTP 服务器（端口从 config.yaml 中读取，所有逻辑在 router 包中）
	if err := router.StartServer(r, cfg); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
