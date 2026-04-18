## Why

当前所有路由配置都写在 `cmd/server/main.go` 文件中(第 67-95 行),随着项目增长,main.go 会变得越来越臃肿,职责不清晰。将路由配置抽离到独立文件可以提升代码的可维护性和可读性,符合单一职责原则。

## What Changes

- 新建 `internal/router/router.go` 文件,专门负责路由配置
- 将 main.go 中的路由注册逻辑(中间件、健康检查、API 路由组、重定向路由)迁移到 router 包
- 提供 `NewRouter` 函数,接收所需的 handler 和配置参数,返回配置好的 `*gin.Engine`
- main.go 中调用 `NewRouter` 替代原有的路由配置代码,保持 main 函数专注于初始化和启动流程

## Capabilities

### New Capabilities
- `router-extraction`: 路由配置独立封装,提供清晰的路由注册接口

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- **Affected Files**: 
  - `cmd/server/main.go`: 简化路由配置部分
  - `internal/router/router.go`: 新建文件,包含路由配置逻辑
- **Dependencies**: 无新增依赖
- **Breaking Changes**: 无,仅内部重构,不影响外部 API
