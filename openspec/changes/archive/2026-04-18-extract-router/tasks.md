## 1. 创建 Router 包

- [x] 1.1 创建 `internal/router/router.go` 文件
- [x] 1.2 定义 `NewRouter` 函数签名,接收所有必需的依赖参数

## 2. 迁移路由配置逻辑

- [x] 2.1 在 `NewRouter` 中初始化 Gin 引擎并设置模式
- [x] 2.2 迁移全局中间件注册(Recovery、Logger、RateLimit)
- [x] 2.3 迁移健康检查端点 `GET /health`
- [x] 2.4 迁移 API 路由组 `/api/v1` 和创建短链端点
- [x] 2.5 迁移短链重定向路由 `GET /:code`

## 3. 更新 Main 函数

- [x] 3.1 在 `main.go` 中导入 `internal/router` 包
- [x] 3.2 替换原有路由配置代码为 `router.NewRouter()` 调用
- [x] 3.3 验证编译通过,无语法错误

## 4. 测试验证

- [x] 4.1 运行 `make build` 确保项目可以正常编译
- [x] 4.2 启动服务并测试健康检查端点 `/health`
- [x] 4.3 测试创建短链 API `POST /api/v1/links`
- [x] 4.4 测试短链重定向功能 `GET /:code`
- [x] 4.5 运行现有测试用例确保无回归
