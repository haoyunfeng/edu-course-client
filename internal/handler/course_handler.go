package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/haoyunfeng/edu-course-client/internal/service"
	coursepb "github.com/haoyunfeng/edu-course-proto/pb"
)

type CourseHandler struct {
	courseService *service.CourseService
}

func NewCourseHandler(courseService *service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// GetCourse 获取课程详情
// @Summary 获取课程详情
// @Description 根据课程ID获取课程详细信息
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "课程ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses/{id} [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseIDStr := c.Param("id")
	if courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course id is required"})
		return
	}

	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	resp, err := h.courseService.GetCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListCourses 获取课程列表
// @Summary 获取课程列表
// @Description 分页获取课程列表
// @Tags courses
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses [get]
func (h *CourseHandler) ListCourses(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "10"), 10, 32)

	resp, err := h.courseService.ListCourses(c.Request.Context(), int32(page), int32(pageSize))
	if err != nil {
		// 记录详细错误信息
		log.Printf("ListCourses error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"message": "Failed to get courses list",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateCourse 创建课程
// @Summary 创建课程
// @Description 创建新课程
// @Tags courses
// @Accept json
// @Produce json
// @Param course body map[string]interface{} true "课程信息"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req coursepb.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.courseService.CreateCourse(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateCourse 更新课程
// @Summary 更新课程
// @Description 更新课程信息
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "课程ID"
// @Param course body map[string]interface{} true "课程信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	courseIDStr := c.Param("id")
	if courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course id is required"})
		return
	}

	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	var req coursepb.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Id = courseID
	resp, err := h.courseService.UpdateCourse(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteCourse 删除课程
// @Summary 删除课程
// @Description 根据课程ID删除课程
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "课程ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseIDStr := c.Param("id")
	if courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course id is required"})
		return
	}

	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	resp, err := h.courseService.DeleteCourse(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
