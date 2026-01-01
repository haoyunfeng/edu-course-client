package registry

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"go-micro.dev/v4/registry"
)

// NewConsulRegistry 创建 Consul 注册中心
func NewConsulRegistry(address string) (registry.Registry, error) {
	// 解析地址（移除 http:// 前缀）
	addr := address
	if strings.HasPrefix(addr, "http://") {
		addr = strings.TrimPrefix(addr, "http://")
	}
	if strings.HasPrefix(addr, "https://") {
		addr = strings.TrimPrefix(addr, "https://")
	}
	// 移除末尾的斜杠
	addr = strings.TrimSuffix(addr, "/")

	// 创建 Consul 客户端配置
	config := api.DefaultConfig()
	config.Address = addr
	config.Scheme = "http"

	// 创建 Consul 客户端
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	// 返回 Consul 注册中心
	return &consulRegistry{
		client: client,
		config: config,
	}, nil
}

type consulRegistry struct {
	client *api.Client
	config *api.Config
}

func (r *consulRegistry) Init(opts ...registry.Option) error {
	return nil
}

func (r *consulRegistry) Options() registry.Options {
	return registry.Options{}
}

func (r *consulRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	if s == nil || len(s.Nodes) == 0 {
		return fmt.Errorf("service or nodes is nil")
	}

	// 为每个节点注册服务
	for _, node := range s.Nodes {
		// 将 registry.Service 转换为 Consul 服务注册
		// 将 Metadata map 转换为 Tags slice
		tags := make([]string, 0, len(s.Metadata))
		for k, v := range s.Metadata {
			tags = append(tags, fmt.Sprintf("%s:%s", k, v))
		}

		registration := &api.AgentServiceRegistration{
			ID:      node.Id,
			Name:    s.Name,
			Tags:    tags,
			Address: node.Address,
		}

		// 从地址中提取端口（格式：host:port）
		if idx := strings.LastIndex(node.Address, ":"); idx > 0 {
			portStr := node.Address[idx+1:]
			var port int
			fmt.Sscanf(portStr, "%d", &port)
			registration.Port = port
			registration.Address = node.Address[:idx]
		}

		// 添加健康检查
		registration.Check = &api.AgentServiceCheck{
			TTL:                            "30s",
			DeregisterCriticalServiceAfter: "90s",
		}

		// 注册服务
		err := r.client.Agent().ServiceRegister(registration)
		if err != nil {
			return fmt.Errorf("failed to register service: %w", err)
		}
	}

	return nil
}

func (r *consulRegistry) Deregister(s *registry.Service, opts ...registry.DeregisterOption) error {
	if s == nil || len(s.Nodes) == 0 {
		return fmt.Errorf("service or nodes is nil")
	}

	// 注销所有节点
	for _, node := range s.Nodes {
		err := r.client.Agent().ServiceDeregister(node.Id)
		if err != nil {
			return fmt.Errorf("failed to deregister service: %w", err)
		}
	}

	return nil
}

func (r *consulRegistry) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	services, _, err := r.client.Health().Service(name, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	// 按服务名称分组
	serviceMap := make(map[string]*registry.Service)
	for _, s := range services {
		svcName := s.Service.Service
		if _, exists := serviceMap[svcName]; !exists {
			// 将 Tags slice 转换为 Metadata map
			metadata := make(map[string]string)
			for _, tag := range s.Service.Tags {
				if idx := strings.Index(tag, ":"); idx > 0 {
					metadata[tag[:idx]] = tag[idx+1:]
				} else {
					metadata[tag] = ""
				}
			}
			// 添加 Meta 信息
			for k, v := range s.Service.Meta {
				metadata[k] = v
			}

			serviceMap[svcName] = &registry.Service{
				Name:     svcName,
				Version:  s.Service.Meta["version"],
				Metadata: metadata,
				Nodes:    []*registry.Node{},
			}
		}

		// 添加节点
		node := &registry.Node{
			Id:       s.Service.ID,
			Address:  fmt.Sprintf("%s:%d", s.Service.Address, s.Service.Port),
			Metadata: s.Service.Meta,
		}
		serviceMap[svcName].Nodes = append(serviceMap[svcName].Nodes, node)
	}

	result := make([]*registry.Service, 0, len(serviceMap))
	for _, s := range serviceMap {
		result = append(result, s)
	}

	return result, nil
}

func (r *consulRegistry) ListServices(opts ...registry.ListOption) ([]*registry.Service, error) {
	services, err := r.client.Agent().Services()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	// 按服务名称分组
	serviceMap := make(map[string]*registry.Service)
	for _, s := range services {
		svcName := s.Service
		if _, exists := serviceMap[svcName]; !exists {
			// 将 Tags slice 转换为 Metadata map
			metadata := make(map[string]string)
			for _, tag := range s.Tags {
				if idx := strings.Index(tag, ":"); idx > 0 {
					metadata[tag[:idx]] = tag[idx+1:]
				} else {
					metadata[tag] = ""
				}
			}
			// 添加 Meta 信息
			for k, v := range s.Meta {
				metadata[k] = v
			}

			serviceMap[svcName] = &registry.Service{
				Name:     svcName,
				Metadata: metadata,
				Nodes:    []*registry.Node{},
			}
		}

		// 添加节点
		node := &registry.Node{
			Id:       s.ID,
			Address:  fmt.Sprintf("%s:%d", s.Address, s.Port),
			Metadata: s.Meta,
		}
		serviceMap[svcName].Nodes = append(serviceMap[svcName].Nodes, node)
	}

	result := make([]*registry.Service, 0, len(serviceMap))
	for _, s := range serviceMap {
		result = append(result, s)
	}

	return result, nil
}

func (r *consulRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return &consulWatcher{
		registry: r,
		stop:     make(chan bool),
	}, nil
}

func (r *consulRegistry) String() string {
	return "consul"
}

type consulWatcher struct {
	registry *consulRegistry
	stop     chan bool
}

func (w *consulWatcher) Next() (*registry.Result, error) {
	// 简单的实现，实际应该使用 Consul 的 Watch API
	select {
	case <-w.stop:
		return nil, fmt.Errorf("watcher stopped")
	case <-time.After(5 * time.Second):
		return &registry.Result{
			Action: "update",
		}, nil
	}
}

func (w *consulWatcher) Stop() {
	close(w.stop)
}
