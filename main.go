package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	consulRegistry "github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	// TODO: 根据实际的 proto 包路径导入服务定义
	// 示例:
	// import course "github.com/haoyunfeng/edu-course-proto/gen/go/course"
	// 请根据实际的 proto 文件调整导入路径
)

var (
	consulAddr  = flag.String("consul", "", "Consul address (default: http://47.113.220.106:8500)")
	serviceName = flag.String("service", "", "Service name to call (for discovery mode)")
	list        = flag.Bool("list", false, "List all available services")
	httpPort    = flag.String("port", "", "HTTP server port (default: 8080)")
	httpMode    = flag.Bool("http", true, "Run HTTP server (default: true)")
)

func main() {
	flag.Parse()

	// 加载配置
	config := LoadConfig()
	if *consulAddr != "" {
		config.ConsulAddr = *consulAddr
	}
	if *serviceName != "" {
		config.ServiceName = *serviceName
	}
	if *httpPort != "" {
		config.HTTPPort = *httpPort
	}

	// 创建 Consul 注册中心
	reg := consulRegistry.NewRegistry(
		registry.Addrs(config.ConsulAddr),
	)

	// 创建微服务客户端
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("edu-course-client"),
	)

	// 初始化服务
	service.Init()

	ctx := context.Background()

	// 如果指定了 list 参数，列出所有服务
	if *list {
		listAllServices(reg)
		return
	}

	// 如果指定了 -service 参数但没有启用 HTTP 模式，使用旧的监控模式
	if config.ServiceName != "" && !*httpMode {
		runDiscoveryMode(ctx, reg, config.ServiceName)
		return
	}

	// 默认启动 HTTP 服务器
	runHTTPServer(ctx, service, config, reg)
}

// runDiscoveryMode 运行服务发现模式（旧模式，用于监控服务）
func runDiscoveryMode(ctx context.Context, reg registry.Registry, serviceName string) {
	// 从 Consul 中查找服务
	services, err := reg.GetService(serviceName)
	if err != nil {
		log.Fatalf("Failed to get service %s from registry: %v", serviceName, err)
	}

	if len(services) == 0 {
		log.Fatalf("Service %s not found in registry", serviceName)
	}

	log.Printf("Found %d instance(s) of service %s:", len(services), serviceName)
	for _, svc := range services {
		log.Printf("  Service: %s, Version: %s", svc.Name, svc.Version)
		for _, node := range svc.Nodes {
			log.Printf("    Node: %s, Address: %s", node.Id, node.Address)
			if len(node.Metadata) > 0 {
				log.Printf("      Metadata: %v", node.Metadata)
			}
		}
	}

	// 定期检查服务可用性
	log.Println("Monitoring service availability... (Press Ctrl+C to exit)")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			services, err := reg.GetService(serviceName)
			if err != nil {
				log.Printf("Error getting service: %v", err)
				continue
			}
			log.Printf("Service %s has %d instance(s) available", serviceName, len(services))
		case <-ctx.Done():
			return
		}
	}
}

// runHTTPServer 运行 HTTP 服务器
func runHTTPServer(ctx context.Context, service micro.Service, config *Config, reg registry.Registry) {
	// 检查目标服务是否可用
	services, err := reg.GetService(config.ServiceAddr)
	if err != nil {
		log.Printf("Warning: Failed to get service %s from registry: %v", config.ServiceAddr, err)
		log.Printf("HTTP server will start anyway, but service calls may fail")
	} else if len(services) == 0 {
		log.Printf("Warning: Service %s not found in registry", config.ServiceAddr)
		log.Printf("HTTP server will start anyway, but service calls may fail")
	} else {
		log.Printf("Service %s is available with %d instance(s)", config.ServiceAddr, len(services))
	}

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 路由器
	r := gin.Default()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 创建处理器并设置路由
	handler := NewHandler(service, config.ServiceAddr)
	handler.SetupRoutes(r)

	// 启动 HTTP 服务器
	addr := fmt.Sprintf(":%s", config.HTTPPort)
	log.Printf("Starting HTTP server on %s", addr)
	log.Printf("API endpoints available at http://localhost%s/api/v1", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// listAllServices 列出 Consul 中所有可用的服务
func listAllServices(reg registry.Registry) {
	services, err := reg.ListServices()
	if err != nil {
		log.Fatalf("Failed to list services: %v", err)
	}

	if len(services) == 0 {
		log.Println("No services found in registry")
		return
	}

	log.Printf("Found %d service(s) in registry:\n", len(services))
	for _, svc := range services {
		log.Printf("  - %s (Version: %s)", svc.Name, svc.Version)
	}
}

// callService 调用服务的示例函数（需要根据实际的 proto 定义调整）
func callService(ctx context.Context, service micro.Service, serviceName string) error {
	// 示例：如何创建一个客户端并调用服务
	// 请根据实际的 proto 定义修改此函数

	// 1. 导入 proto 包（在文件顶部）
	// import course "github.com/haoyunfeng/edu-course-proto/gen/go/course"

	// 2. 创建客户端（需要根据实际的 proto 定义）
	// client := course.NewCourseService(serviceName, service.Client())

	// 3. 准备请求（需要根据实际的 proto 定义）
	// req := &course.GetCourseRequest{
	//     Id: "123",
	// }

	// 4. 调用服务
	// rsp, err := client.GetCourse(ctx, req)
	// if err != nil {
	//     return fmt.Errorf("failed to call service: %w", err)
	// }

	// 5. 处理响应
	// log.Printf("Response: %+v", rsp)

	log.Println("callService function needs to be implemented based on actual proto definitions")
	return nil
}
