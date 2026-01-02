# 路由系统说明

## 概述

本路由系统采用集中式注册方式，所有路由注册逻辑统一在 `router/register.go` 中管理，使得添加新路由变得简单和灵活。

## 架构

### 1. Route 接口

所有路由都需要实现 `Route` 接口：

```go
type Route interface {
    Register(group *gin.RouterGroup)
}
```

### 2. Router 管理器

`Router` 负责管理所有路由的注册和设置：

- `NewRouter(engine *gin.Engine)`: 创建路由管理器
- `Register(route Route)`: 注册路由
- `Setup()`: 设置所有路由

### 3. Services 服务集合

`Services` 结构体用于统一管理所有服务依赖，方便路由注册时使用。

## 使用方式

### 在 main.go 中使用（极简）

```go
// 统一设置所有路由和服务（只需一行代码）
services, err := router.SetupRoutes(r, cfg)
if err != nil {
    log.Fatalf("Failed to setup routes: %v", err)
}
defer services.Close()
```

**main.go 无需任何服务初始化代码！** 所有服务初始化和路由注册都在 `router.SetupRoutes()` 中完成。

### 添加新路由

**只需在 `router/register.go` 中添加即可，无需修改 main.go**：

1. **在 Services 结构体中添加服务字段**：

```go
type Services struct {
    CourseService *service.CourseService
    UserService   *service.UserService  // 新增
}
```

2. **在 initServices 函数中初始化服务**：

```go
func initServices(cfg *config.Config) (*Services, error) {
    services := &Services{}

    // 初始化课程服务
    courseService, err := service.NewCourseService(cfg)
    if err != nil {
        return nil, err
    }
    services.CourseService = courseService

    // 初始化用户服务（新增）
    userService, err := service.NewUserService(cfg)
    if err != nil {
        return nil, err
    }
    services.UserService = userService

    return services, nil
}
```

3. **创建路由注册函数**：

```go
// registerUserRoutes 注册用户相关路由
func registerUserRoutes(router *Router, userService *service.UserService) {
    userHandler := handler.NewUserHandler(userService)
    
    userRoute := &userRoute{
        handler: userHandler,
    }
    router.Register(userRoute)
}

// userRoute 用户路由实现
type userRoute struct {
    handler *handler.UserHandler
}

func (r *userRoute) Register(group *gin.RouterGroup) {
    users := group.Group("/users")
    {
        users.GET("/:id", r.handler.GetUser)
        users.GET("", r.handler.ListUsers)
        users.POST("", r.handler.CreateUser)
        users.PUT("/:id", r.handler.UpdateUser)
        users.DELETE("/:id", r.handler.DeleteUser)
    }
}
```

4. **在 SetupRoutes 函数中注册路由**：

```go
func SetupRoutes(r *gin.Engine, cfg *config.Config) (*Services, error) {
    // 初始化所有服务
    services, err := initServices(cfg)
    if err != nil {
        return nil, err
    }

    router := NewRouter(r)

    // 注册课程路由
    if services.CourseService != nil {
        registerCourseRoutes(router, services.CourseService)
    }

    // 注册用户路由（新增）
    if services.UserService != nil {
        registerUserRoutes(router, services.UserService)
    }

    router.Setup()
    return services, nil
}
```

5. **在 Services.Close() 中添加服务关闭逻辑**：

```go
func (s *Services) Close() {
    if s.CourseService != nil {
        s.CourseService.Close()
    }
    // 关闭用户服务（新增）
    if s.UserService != nil {
        s.UserService.Close()
    }
}
```

**完成！main.go 无需任何修改！**

## 优势

1. **集中管理**: 所有路由注册逻辑集中在 `register.go` 中
2. **main.go 简洁**: main.go 只需调用一个函数，无需逐个注册
3. **易于扩展**: 添加新路由只需在 `register.go` 中添加，无需修改 main.go
4. **类型安全**: 通过接口确保路由实现正确
5. **避免循环依赖**: 路由实现放在 router 包中，避免 handler 包依赖 router 包

## 示例

参考 `internal/router/register.go` 中的 `registerCourseRoutes` 函数查看完整的实现示例。
