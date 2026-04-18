## Context

当前 `cmd/server/main.go` 文件中包含了完整的应用初始化流程:配置加载、数据库初始化、Redis 初始化、依赖注入、路由配置和服务器启动。其中路由配置部分(第 67-95 行)占了约 30 行代码,包含中间件注册、健康检查端点、API 路由组和重定向路由。

随着项目发展,可能会添加更多路由和中间件,main.go 会越来越长,不利于维护和测试。

## Goals / Non-Goals

**Goals:**
- 将路由配置逻辑抽离到独立的 `internal/router` 包
- 保持现有的路由结构和中间件配置不变
- 提供清晰的函数接口,便于测试和扩展
- 简化 main.go,使其专注于初始化和生命周期管理

**Non-Goals:**
- 不改变现有的路由路径或中间件逻辑
- 不引入新的依赖或框架
- 不做其他重构(如 handler 改造、配置重构等)

## Decisions

### 1. 函数签名设计

**决策**: `NewRouter` 函数接收 Repository 和配置参数,在内部初始化 Handler

```go
func NewRouter(
    linkRepo repository.LinkRepositoryInterface,
    shortLinkConfig *config.ShortLinkConfig,
    rateLimitConfig *config.RateLimitConfig,
    logger *zap.Logger,
) *gin.Engine
```

**理由**: 
- Handler 的初始化属于应用组装逻辑,应该集中在 Router 层
- 减少 main.go 的职责,只需传递基础设施层依赖(Repository、Config、Logger)
- 符合常见的分层架构模式:Router 负责组装 Handler → Service → Repository

**备选方案**: 
- 方案 A: 在 main.go 中初始化 Handler 后传入 → 会导致 main.go 臃肿
- 方案 B: 使用全局变量 → 不利于测试和维护

### 2. 包位置

**决策**: 放在 `internal/router/router.go`

**理由**:
- `internal` 包表示内部实现,不对外暴露
- 与现有的 `handler`、`service`、`middleware` 等包保持一致的层次结构
- 符合 Go 项目的常见组织方式

### 4. Router 内部依赖组装

**决策**: Router 函数内部按以下顺序初始化依赖链

```go
// 1. 初始化 Service
linkService := service.NewLinkService(linkRepo, shortLinkConfig)

// 2. 初始化 Handlers
linkHandler := handler.NewLinkHandler(linkService)
redirectHandler := handler.NewRedirectHandler(linkService)

// 3. 配置路由和中间件
// ...
```

**理由**: 
- 保持依赖注入的清晰流向: Repository → Service → Handler
- Router 作为组装层,负责将所有组件连接在一起
- main.go 只需关注基础设施初始化(DB、Redis、Config)

### 3. 中间件注册方式

**决策**: 保持现有的中间件注册顺序和方式不变

```go
router.Use(middleware.Recovery(logger))
router.Use(middleware.Logger(logger))
if rateLimitEnabled {
    router.Use(middleware.RateLimit(rate, burst))
}
```

**理由**: 中间件顺序很重要(Recovery 必须最先),保持现有逻辑避免引入 bug

## Risks / Trade-offs

### [低风险] 参数列表较长
**影响**: `NewRouter` 函数有 6 个参数,可能不够简洁  
**缓解**: 可以考虑使用配置结构体,但对于当前场景,显式参数更清晰

### [无风险] 路由逻辑分散
**影响**: 有人可能认为路由应该集中在 main.go 更直观  
**缓解**: 这是 Go 项目的常见实践,大多数 Web 框架示例都采用这种方式

## Migration Plan

1. 创建 `internal/router/router.go` 文件
2. 将路由配置代码迁移到新文件
3. 修改 `main.go` 调用 `router.NewRouter()`
4. 运行测试验证功能正常
5. 提交代码

**回滚策略**: 如果出现问题,可以直接将路由代码移回 main.go,因为接口完全兼容

## Open Questions

无
