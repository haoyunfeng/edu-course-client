package router

import (
	"github.com/gin-gonic/gin"
	"github.com/haoyunfeng/edu-course-client/internal/config"
	"github.com/haoyunfeng/edu-course-client/internal/handler"
	"github.com/haoyunfeng/edu-course-client/internal/service"
)

// SetupRoutes 设置所有路由并启动服务器
// 统一管理路由、服务初始化和服务器启动，main.go 只需调用此函数即可
func SetupRoutes(r *gin.Engine, cfg *config.Config) (*Services, error) {
	// 初始化所有服务
	services, err := initServices(cfg)
	if err != nil {
		return nil, err
	}

	// 创建路由管理器
	router := NewRouter(r)

	// 注册课程路由
	if services.CourseService != nil {
		registerCourseRoutes(router, services.CourseService)
	}

	// 在这里添加更多路由注册
	// 例如：
	// if services.UserService != nil {
	//     registerUserRoutes(router, services.UserService)
	// }

	// 设置所有路由
	router.Setup()

	return services, nil
}

// StartServer 启动 HTTP 服务器
// 从 config.yaml 中读取端口并启动服务器
// 如果配置文件中未指定端口，使用默认端口 :8083
func StartServer(r *gin.Engine, cfg *config.Config) error {
	addr := cfg.HTTP.Address
	if addr == "" {
		addr = ":8083" // 默认端口（如果配置文件中未指定）
	}
	return r.Run(addr)
}

// initServices 初始化所有服务
// 在这里统一管理所有服务的初始化，新增服务只需在此函数中添加即可
func initServices(cfg *config.Config) (*Services, error) {
	services := &Services{}

	// 初始化课程服务
	courseService, err := service.NewCourseService(cfg)
	if err != nil {
		return nil, err
	}
	services.CourseService = courseService

	// 在这里添加更多服务初始化
	// 例如：
	// userService, err := service.NewUserService(cfg)
	// if err != nil {
	//     return nil, err
	// }
	// services.UserService = userService

	return services, nil
}

// registerCourseRoutes 注册课程相关路由
func registerCourseRoutes(router *Router, courseService *service.CourseService) {
	courseHandler := handler.NewCourseHandler(courseService)

	// 创建课程路由
	courseRoute := &courseRoute{
		handler: courseHandler,
	}
	router.Register(courseRoute)
}

// courseRoute 课程路由实现
type courseRoute struct {
	handler *handler.CourseHandler
}

func (r *courseRoute) Register(group *gin.RouterGroup) {
	courses := group.Group("/courses")
	{
		courses.GET("/:id", r.handler.GetCourse)
		courses.GET("", r.handler.ListCourses)
		courses.POST("", r.handler.CreateCourse)
		courses.PUT("/:id", r.handler.UpdateCourse)
		courses.DELETE("/:id", r.handler.DeleteCourse)
	}
}

// Services 服务依赖集合
// 用于统一管理所有服务，方便路由注册时使用
type Services struct {
	CourseService *service.CourseService
	// 在这里添加更多服务
	// UserService    *service.UserService
	// OrderService   *service.OrderService
}

// Close 关闭所有服务
func (s *Services) Close() {
	if s.CourseService != nil {
		s.CourseService.Close()
	}
	// 在这里添加更多服务的关闭逻辑
	// if s.UserService != nil {
	//     s.UserService.Close()
	// }
}
