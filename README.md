# edu-course-client

基于 go-micro 和 Consul 的微服务客户端应用，用于发现和调用注册在 Consul 上的微服务。

## 功能特性

- ✅ 使用 Consul 作为服务发现注册中心
- ✅ 支持通过命令行参数或环境变量配置
- ✅ 服务发现和监控
- ✅ 列出所有可用服务
- ✅ 实时监控服务可用性

## 前置要求

- Go 1.16 或更高版本
- Consul 服务运行在 `http://47.113.220.106:8500`
- 访问 `github.com/haoyunfeng/edu-course-proto/gen/go@v1.0.0` 的权限

## 安装

```bash
# 克隆项目
git clone git@github.com:haoyunfeng/edu-course-client.git
cd edu-course-client

# 下载依赖
go mod download

# 编译
go build -o edu-course-client
```

## 使用方法

### 启动 HTTP 服务（推荐）

```bash
# 默认启动 HTTP 服务，监听 8080 端口
./edu-course-client

# 指定端口
./edu-course-client -port 3000

# 指定 Consul 地址和端口
./edu-course-client -consul http://47.113.220.106:8500 -port 8080
```

### 服务发现模式

```bash
# 查看所有可用服务
./edu-course-client -list

# 查找并监控特定服务（不使用 HTTP 模式）
./edu-course-client -http=false -service <服务名称>
```

### 使用环境变量

```bash
# 设置环境变量
export CONSUL_ADDR=http://47.113.220.106:8500
export HTTP_PORT=8080
export SERVICE_ADDR=edu-course

# 运行 HTTP 服务
./edu-course-client
```

### 命令行参数

- `-consul`: Consul 服务地址（默认: `http://47.113.220.106:8500`）
- `-port`: HTTP 服务器端口（默认: `8080`）
- `-http`: 是否启动 HTTP 服务器（默认: `true`）
- `-service`: 要查找的服务名称（用于服务发现模式）
- `-list`: 列出所有在 Consul 中注册的服务

### HTTP API 接口

服务启动后，可以通过以下 API 接口访问：

#### 健康检查
```bash
GET /api/v1/health
```

#### 课程相关接口
```bash
# 获取课程列表
GET /api/v1/courses

# 获取课程详情
GET /api/v1/courses/:id

# 创建课程
POST /api/v1/courses
Content-Type: application/json

# 更新课程
PUT /api/v1/courses/:id
Content-Type: application/json

# 删除课程
DELETE /api/v1/courses/:id
```

### API 示例

```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# 获取课程详情
curl http://localhost:8080/api/v1/courses/123

# 获取课程列表
curl http://localhost:8080/api/v1/courses
```

## 配置说明

应用使用 `config.yaml` 文件进行配置。配置文件优先级：命令行参数 > 环境变量 > 配置文件 > 默认值

### 配置文件 (config.yaml)

创建 `config.yaml` 文件（项目根目录已包含示例）：

```yaml
# Consul 配置
consul:
  addr: "http://47.113.220.106:8500"

# HTTP 服务器配置
server:
  port: "8080"

# 微服务配置
service:
  # 要调用的微服务名称
  addr: "edu-course"
  # 服务发现模式使用的服务名称（可选）
  name: ""
```

### 配置项说明

#### Consul 地址

- **配置文件**: `consul.addr`
- **环境变量**: `CONSUL_ADDR`
- **命令行参数**: `-consul <地址>`
- **默认值**: `http://47.113.220.106:8500`

#### HTTP 端口

- **配置文件**: `server.port`
- **环境变量**: `HTTP_PORT`
- **命令行参数**: `-port <端口>`
- **默认值**: `8080`

#### 微服务名称

要调用的微服务名称，应该与在 Consul 中注册的服务名称一致。

- **配置文件**: `service.addr`
- **环境变量**: `SERVICE_ADDR`
- **默认值**: `edu-course`

#### 服务发现模式服务名

用于服务发现模式的服务名称（可选）。

- **配置文件**: `service.name`
- **环境变量**: `SERVICE_NAME`
- **命令行参数**: `-service <服务名称>`

### 配置优先级

1. **命令行参数** - 最高优先级
2. **环境变量** - 次优先级
3. **config.yaml** - 配置文件
4. **默认值** - 最低优先级

## 调用服务

目前代码提供了服务发现和监控功能。要实际调用服务，需要根据 `github.com/haoyunfeng/edu-course-proto/gen/go@v1.0.0` 中定义的 proto 文件来修改代码。

### 示例：如何调用服务

1. **导入 proto 包**

   在 `main.go` 文件顶部添加导入：

   ```go
   import course "github.com/haoyunfeng/edu-course-proto/gen/go/course"
   ```

   （请根据实际的 proto 包路径调整）

2. **创建客户端并调用服务**

   在 `main()` 函数或 `callService()` 函数中添加代码：

   ```go
   // 创建客户端
   client := course.NewCourseService(config.ServiceName, service.Client())
   
   // 准备请求
   req := &course.GetCourseRequest{
       Id: "123",
   }
   
   // 调用服务
   rsp, err := client.GetCourse(ctx, req)
   if err != nil {
       log.Fatalf("Failed to call service: %v", err)
   }
   
   // 处理响应
   log.Printf("Response: %+v", rsp)
   ```

   （请根据实际的 proto 定义调整服务名、请求和响应类型）

## 项目结构

```
edu-course-client/
├── main.go          # 主程序文件
├── config.go        # 配置管理
├── config.yaml      # 配置文件
├── handler.go       # HTTP 请求处理器
├── go.mod           # Go 模块定义
├── go.sum           # 依赖校验和
└── README.md        # 项目说明文档
```

## 依赖项

- `go-micro.dev/v4`: go-micro 微服务框架
- `github.com/go-micro/plugins/v4/registry/consul`: Consul 注册中心插件
- `github.com/gin-gonic/gin`: Gin HTTP 框架
- `gopkg.in/yaml.v3`: YAML 配置文件解析
- `github.com/haoyunfeng/edu-course-proto/gen/go@v1.0.0`: Proto 定义包

## 开发

```bash
# 运行 HTTP 服务（开发模式）
go run main.go config.go handler.go

# 指定端口运行
go run main.go config.go handler.go -port 3000

# 运行测试
go test ./...

# 代码格式化
go fmt ./...

# 代码检查
go vet ./...

# 编译
go build -o edu-course-client
```

## 注意事项

1. 确保 Consul 服务正在运行并可以访问
2. 确保目标微服务已注册到 Consul
3. 根据实际的 proto 定义修改服务调用代码
4. 生产环境建议添加适当的错误处理和重试机制

## 许可证

[添加许可证信息]
