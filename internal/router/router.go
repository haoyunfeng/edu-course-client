package router

import (
	"github.com/gin-gonic/gin"
)

// Route 路由定义接口
type Route interface {
	// Register 注册路由到路由组
	Register(group *gin.RouterGroup)
}

// Router 路由管理器
type Router struct {
	engine *gin.Engine
	routes []Route
}

// NewRouter 创建新的路由管理器
func NewRouter(engine *gin.Engine) *Router {
	return &Router{
		engine: engine,
		routes: make([]Route, 0),
	}
}

// Register 注册路由
func (r *Router) Register(route Route) {
	r.routes = append(r.routes, route)
}

// Setup 设置所有路由
func (r *Router) Setup() {
	// 健康检查
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由组
	api := r.engine.Group("/api/v1")
	{
		// 注册所有路由
		for _, route := range r.routes {
			route.Register(api)
		}
	}
}
