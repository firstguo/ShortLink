## ADDED Requirements

### Requirement: Router 初始化函数

系统 SHALL 提供 `NewRouter` 函数,接收 Repository 和配置参数,在内部初始化完整的依赖链(Service 和 Handler)并返回配置好的 `*gin.Engine` 实例。

#### Scenario: 创建完整配置的路由器

- **WHEN** 调用 `NewRouter` 函数,传入 linkRepo、shortLinkConfig、rateLimitConfig 和 logger
- **THEN** 返回一个配置了所有中间件和路由的 `*gin.Engine` 实例,内部自动初始化 Service 和 Handler

#### Scenario: 禁用限流时不注册限流中间件

- **WHEN** 调用 `NewRouter` 时 `rateLimitConfig.Enabled` 为 false
- **THEN** 返回的路由器不包含 RateLimit 中间件

### Requirement: Router 内部依赖组装

Router SHALL 在内部按以下顺序初始化依赖链:
1. 使用 linkRepo 和 shortLinkConfig 初始化 LinkService
2. 使用 LinkService 初始化 LinkHandler 和 RedirectHandler
3. 配置路由和中间件

#### Scenario: 依赖链正确组装

- **WHEN** 调用 `NewRouter`
- **THEN** 内部创建 LinkService → LinkHandler 和 RedirectHandler 的依赖链

### Requirement: 中间件注册顺序

路由器 SHALL 按照以下顺序注册全局中间件:
1. Recovery 中间件(必须最先注册)
2. Logger 中间件
3. RateLimit 中间件(可选,仅在启用时注册)

#### Scenario: 中间件按正确顺序注册

- **WHEN** 启用所有中间件并调用 `NewRouter`
- **THEN** 中间件的执行顺序为: Recovery → Logger → RateLimit

### Requirement: 健康检查端点

路由器 SHALL 注册 `GET /health` 端点,返回服务健康状态和当前时间。

#### Scenario: 健康检查返回成功状态

- **WHEN** 发送 GET 请求到 `/health`
- **THEN** 返回 HTTP 200 状态码,响应体包含 `status: "ok"` 和 `time` 字段

### Requirement: API 路由组

路由器 SHALL 注册 `/api/v1` 路由组,包含创建短链的 POST 端点。

#### Scenario: 创建短链端点注册

- **WHEN** 调用 `NewRouter`
- **THEN** `POST /api/v1/links` 端点被正确注册并路由到 LinkHandler.CreateLink

### Requirement: 短链重定向路由

路由器 SHALL 在根路径注册 `GET /:code` 端点用于短链重定向。

#### Scenario: 重定向端点注册

- **WHEN** 调用 `NewRouter`
- **THEN** `GET /:code` 端点被正确注册并路由到 RedirectHandler.Redirect
