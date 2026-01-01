package service

import (
	"context"
	"fmt"

	grpcClient "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/haoyunfeng/edu-course-client/internal/config"
	"github.com/haoyunfeng/edu-course-client/internal/registry"
	coursepb "github.com/haoyunfeng/edu-course-proto/pb"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	microRegistry "go-micro.dev/v4/registry"
)

// CourseServiceClient 课程服务客户端接口
// 注意：这里使用通用接口，实际使用时需要根据 edu-course-proto 中的具体定义进行调整
type CourseServiceClient interface {
	GetCourse(ctx context.Context, req interface{}, opts ...client.CallOption) (interface{}, error)
	ListCourses(ctx context.Context, req interface{}, opts ...client.CallOption) (interface{}, error)
	CreateCourse(ctx context.Context, req interface{}, opts ...client.CallOption) (interface{}, error)
	UpdateCourse(ctx context.Context, req interface{}, opts ...client.CallOption) (interface{}, error)
	DeleteCourse(ctx context.Context, req interface{}, opts ...client.CallOption) error
}

type CourseService struct {
	client client.Client
	cfg    *config.Config
}

func NewCourseService(cfg *config.Config) (*CourseService, error) {
	// 根据配置选择注册中心
	var reg microRegistry.Registry
	var err error
	switch cfg.Micro.Registry {
	case "consul":
		// 使用 Consul 作为注册中心
		if cfg.Micro.RegistryAddr == "" {
			return nil, fmt.Errorf("consul registry address is required")
		}
		reg, err = registry.NewConsulRegistry(cfg.Micro.RegistryAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to create consul registry: %w", err)
		}
	case "mdns":
		// 使用 mDNS 作为注册中心（默认）
		fallthrough
	default:
		// 默认使用 mDNS
		reg = nil // go-micro 默认使用 mDNS
	}

	// 创建 micro 服务
	opts := []micro.Option{
		micro.Name("edu.course.client"),
		// 使用 gRPC 客户端（因为服务端使用 gRPC）
		micro.Client(grpcClient.NewClient()),
	}

	// 如果指定了注册中心，则使用它
	if reg != nil {
		opts = append(opts, micro.Registry(reg))
	}

	service := micro.NewService(opts...)

	// 初始化服务
	service.Init()

	return &CourseService{
		client: service.Client(),
		cfg:    cfg,
	}, nil
}

func (s *CourseService) Close() error {
	return nil
}

// GetCourse 获取课程信息
func (s *CourseService) GetCourse(ctx context.Context, courseID int64) (*coursepb.GetCourseResponse, error) {
	req := &coursepb.GetCourseRequest{
		Id: courseID,
	}

	// 使用 go-micro client 调用服务
	// gRPC 方法名格式：/package.Service/Method
	microReq := s.client.NewRequest(
		s.cfg.Micro.CourseServiceName,
		"/course.CourseService/GetCourse",
		req,
	)

	var resp coursepb.GetCourseResponse
	if err := s.client.Call(ctx, microReq, &resp); err != nil {
		return nil, fmt.Errorf("failed to call course service: %w", err)
	}

	return &resp, nil
}

// ListCourses 获取课程列表（使用 GetAllCourses）
func (s *CourseService) ListCourses(ctx context.Context, page, pageSize int32) (*coursepb.GetAllCoursesResponse, error) {
	req := &coursepb.GetAllCoursesRequest{}

	// 使用 go-micro client 调用服务
	// gRPC 方法名格式：/package.Service/Method
	microReq := s.client.NewRequest(
		s.cfg.Micro.CourseServiceName,
		"/course.CourseService/GetAllCourses",
		req,
	)

	var resp coursepb.GetAllCoursesResponse
	if err := s.client.Call(ctx, microReq, &resp); err != nil {
		return nil, fmt.Errorf("failed to call course service %s: %w", s.cfg.Micro.CourseServiceName, err)
	}

	return &resp, nil
}

// CreateCourse 创建课程
func (s *CourseService) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*coursepb.CreateCourseResponse, error) {
	// 使用 go-micro client 调用服务
	// gRPC 方法名格式：/package.Service/Method
	microReq := s.client.NewRequest(
		s.cfg.Micro.CourseServiceName,
		"/course.CourseService/CreateCourse",
		req,
	)

	var resp coursepb.CreateCourseResponse
	if err := s.client.Call(ctx, microReq, &resp); err != nil {
		return nil, fmt.Errorf("failed to call course service: %w", err)
	}

	return &resp, nil
}

// UpdateCourse 更新课程
func (s *CourseService) UpdateCourse(ctx context.Context, req *coursepb.UpdateCourseRequest) (*coursepb.UpdateCourseResponse, error) {
	// 使用 go-micro client 调用服务
	// gRPC 方法名格式：/package.Service/Method
	microReq := s.client.NewRequest(
		s.cfg.Micro.CourseServiceName,
		"/course.CourseService/UpdateCourse",
		req,
	)

	var resp coursepb.UpdateCourseResponse
	if err := s.client.Call(ctx, microReq, &resp); err != nil {
		return nil, fmt.Errorf("failed to call course service: %w", err)
	}

	return &resp, nil
}

// DeleteCourse 删除课程
func (s *CourseService) DeleteCourse(ctx context.Context, courseID int64) (*coursepb.DeleteCourseResponse, error) {
	req := &coursepb.DeleteCourseRequest{
		Id: courseID,
	}

	// 使用 go-micro client 调用服务
	// gRPC 方法名格式：/package.Service/Method
	microReq := s.client.NewRequest(
		s.cfg.Micro.CourseServiceName,
		"/course.CourseService/DeleteCourse",
		req,
	)

	var resp coursepb.DeleteCourseResponse
	if err := s.client.Call(ctx, microReq, &resp); err != nil {
		return nil, fmt.Errorf("failed to call course service: %w", err)
	}

	return &resp, nil
}
