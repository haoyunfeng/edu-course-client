package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haoyunfeng/edu-course-client/internal/config"
	"github.com/haoyunfeng/edu-course-client/internal/handler"
	"github.com/haoyunfeng/edu-course-client/internal/service"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化课程服务客户端
	courseService, err := service.NewCourseService(cfg)
	if err != nil {
		log.Fatalf("Failed to create course service: %v", err)
	}
	defer courseService.Close()

	// 创建 HTTP 处理器
	courseHandler := handler.NewCourseHandler(courseService)

	// 设置 gin 模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 gin 路由
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 课程相关接口
		courses := api.Group("/courses")
		{
			courses.GET("/:id", courseHandler.GetCourse)
			courses.GET("", courseHandler.ListCourses)
			courses.POST("", courseHandler.CreateCourse)
			courses.PUT("/:id", courseHandler.UpdateCourse)
			courses.DELETE("/:id", courseHandler.DeleteCourse)
		}
	}

	// 启动 HTTP 服务器
	addr := cfg.HTTP.Address
	if addr == "" {
		addr = ":8083"
	}

	log.Printf("HTTP server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
