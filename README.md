# Edu Course Client

基于 Go-Micro 和 Gin 的课程服务 HTTP 客户端。

## 功能特性

- 使用 Go-Micro 作为微服务客户端
- 使用 Gin 提供 HTTP RESTful API
- 支持调用课程服务（edu-course-proto）
- 支持服务发现（支持 Consul 和 mDNS）

## 项目结构

```
.
├── main.go                    # 应用入口
├── go.mod                     # Go 模块定义
├── internal/
│   ├── config/               # 配置管理
│   │   └── config.go
│   ├── handler/              # HTTP 处理器
│   │   └── course_handler.go
│   └── service/              # 微服务客户端
│       └── course_service.go
└── README.md
```

## 快速开始

### 1. 安装依赖

项目使用 `github.com/haoyunfeng/edu-course-proto@v1.0.1` 作为 proto 模块。

**Proto 模块 Git 地址：**
- SSH: `git@github.com:haoyunfeng/edu-course-proto.git`
- HTTPS: `https://github.com/haoyunfeng/edu-course-proto.git`

如果遇到 proto 模块校验失败的问题，可以使用以下命令跳过校验：

```bash
# 跳过校验下载依赖
GOSUMDB=off go mod download

# 或者使用直接模式
GOPROXY=direct GOSUMDB=off go mod tidy
```

`go.mod` 文件中已配置 replace 指令指向具体的 commit hash，如果仍有问题，可以手动更新：

```go
replace github.com/haoyunfeng/edu-course-proto => github.com/haoyunfeng/edu-course-proto v1.0.1-0.20260101140148-b2b603ea01c4
```

### 2. 配置

项目支持通过 `config.yaml` 文件或环境变量进行配置。

**方式一：使用 config.yaml（推荐）**

编辑项目根目录下的 `config.yaml` 文件：

```yaml
micro:
  registry: consul                    # 注册中心类型: consul 或 mdns
  registry_addr: "http://47.113.220.106:8500"  # Consul 地址（使用 consul 时必需）
  course_service_name: "edu-course"

http:
  address: ":8083"
```

**方式二：使用环境变量**

如果 `config.yaml` 不存在或配置项为空，将使用环境变量或默认值：

```bash
# 配置文件路径（默认: config.yaml）
export CONFIG_FILE=config.yaml

# 微服务注册中心类型（默认: mdns，可选: consul）
export MICRO_REGISTRY=consul

# 微服务注册中心地址（使用 consul 时必需，格式: http://host:port）
export MICRO_REGISTRY_ADDRESS=http://47.113.220.106:8500

# 课程服务名称（默认: edu-course）
export COURSE_SERVICE_NAME=edu-course

# HTTP 服务器地址（默认: :8083）
export HTTP_ADDRESS=:8083
```

**配置优先级：** `config.yaml` > 环境变量 > 默认值

### 3. 运行服务

```bash
go run main.go
```

## API 接口

### 健康检查

```
GET /health
```

### 课程接口

#### 获取课程详情
```
GET /api/v1/courses/:id
```

#### 获取课程列表
```
GET /api/v1/courses?page=1&pageSize=10
```

#### 创建课程
```
POST /api/v1/courses
Content-Type: application/json

{
  "name": "课程名称",
  "description": "课程描述"
}
```

#### 更新课程
```
PUT /api/v1/courses/:id
Content-Type: application/json

{
  "name": "更新后的课程名称",
  "description": "更新后的课程描述"
}
```

#### 删除课程
```
DELETE /api/v1/courses/:id
```

## 依赖

- [go-micro/v4](https://github.com/go-micro/go-micro) - 微服务框架
- [gin](https://github.com/gin-gonic/gin) - HTTP Web 框架
- [edu-course-proto](https://github.com/haoyunfeng/edu-course-proto) - 课程服务 Proto 定义

## 开发

### 构建

```bash
go build -o bin/edu-course-client main.go
```

### 运行测试

```bash
go test ./...
```

## 许可证

MIT
