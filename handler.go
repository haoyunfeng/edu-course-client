package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
)

// Handler HTTP 请求处理器
type Handler struct {
	service     micro.Service
	serviceName string // 微服务名称
}

// NewHandler 创建新的 HTTP 处理器
func NewHandler(service micro.Service, serviceName string) *Handler {
	return &Handler{
		service:     service,
		serviceName: serviceName,
	}
}

// SetupRoutes 设置路由
func (h *Handler) SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// 健康检查
		api.GET("/health", h.HealthCheck)

		// 课程相关接口
		courses := api.Group("/courses")
		{
			courses.GET("", h.ListCourses)
			courses.GET("/:id", h.GetCourse)
			courses.POST("", h.CreateCourse)
			courses.PUT("/:id", h.UpdateCourse)
			courses.DELETE("/:id", h.DeleteCourse)
		}
	}
}

// HealthCheck 健康检查接口
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "edu-course-client",
	})
}

// GetCourse 获取课程详情
func (h *Handler) GetCourse(c *gin.Context) {
	id := c.Param("id")

	// TODO: 根据实际的 proto 定义调用服务
	// 示例代码（需要根据实际的 proto 定义调整）:
	//
	// 1. 导入 proto 包（在文件顶部）:
	// import course "github.com/haoyunfeng/edu-course-proto/gen/go/course"
	//
	// 2. 创建客户端:
	// client := course.NewCourseService("edu-course", h.service.Client())
	//
	// 3. 准备请求:
	// req := &course.GetCourseRequest{
	//     Id: id,
	// }
	//
	// 4. 调用服务:
	// rsp, err := client.GetCourse(c.Request.Context(), req)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{
	//         "error": "Failed to get course",
	//         "details": err.Error(),
	//     })
	//     return
	// }
	//
	// 5. 返回响应:
	// c.JSON(http.StatusOK, rsp)

	// 临时返回（需要根据实际的 proto 定义替换）
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Course service call - needs proto implementation",
		"note":    "Please implement based on github.com/haoyunfeng/edu-course-proto/gen/go@v1.0.0",
	})
}

// ListCourses 获取课程列表 - 调用 GetAllCourses 方法
// 调用 edu-course 服务的 /course.CourseService/GetAllCourses 端点
func (h *Handler) ListCourses(c *gin.Context) {
	ctx := c.Request.Context()

	// 使用 go-micro 客户端调用服务
	// 服务名: h.serviceName (从配置读取，默认 "edu-course")
	// 端点: "/course.CourseService/GetAllCourses"
	// 请求参数（如果需要分页等参数，可以从 query string 获取）
	req := map[string]interface{}{}

	// 可以从查询参数获取分页信息
	if page := c.Query("page"); page != "" {
		req["page"] = getIntQuery(c, "page", 1)
	}
	if limit := c.Query("limit"); limit != "" {
		req["limit"] = getIntQuery(c, "limit", 10)
	}

	// 创建请求对象
	// NewRequest(service, endpoint, request interface{}, opts ...RequestOption)
	// 端点格式: "/course.CourseService/GetAllCourses"
	request := h.service.Client().NewRequest(
		h.serviceName,                         // 服务名
		"/course.CourseService/GetAllCourses", // 端点名（gRPC 风格路径）
		req,                                   // 请求参数
	)

	// 调用服务
	var rsp interface{}
	if err := h.service.Client().Call(ctx, request, &rsp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get courses",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rsp)
}

// CreateCourse 创建课程
func (h *Handler) CreateCourse(c *gin.Context) {
	// TODO: 根据实际的 proto 定义调用服务
	// var req course.CreateCourseRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	//     c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//     return
	// }
	// client := course.NewCourseService("edu-course", h.service.Client())
	// rsp, err := client.CreateCourse(c.Request.Context(), &req)
	// ...

	// 临时返回
	c.JSON(http.StatusCreated, gin.H{
		"message": "Create course service call - needs proto implementation",
	})
}

// UpdateCourse 更新课程
func (h *Handler) UpdateCourse(c *gin.Context) {
	id := c.Param("id")

	// TODO: 根据实际的 proto 定义调用服务
	// var req course.UpdateCourseRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	//     c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//     return
	// }
	// req.Id = id
	// client := course.NewCourseService("edu-course", h.service.Client())
	// rsp, err := client.UpdateCourse(c.Request.Context(), &req)
	// ...

	// 临时返回
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Update course service call - needs proto implementation",
	})
}

// DeleteCourse 删除课程
func (h *Handler) DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	// TODO: 根据实际的 proto 定义调用服务
	// client := course.NewCourseService("edu-course", h.service.Client())
	// req := &course.DeleteCourseRequest{Id: id}
	// _, err := client.DeleteCourse(c.Request.Context(), req)
	// ...

	// 临时返回
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Delete course service call - needs proto implementation",
	})
}

// getIntQuery 获取整数查询参数
func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}
