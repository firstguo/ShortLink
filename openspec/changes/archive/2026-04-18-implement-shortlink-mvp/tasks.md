## 1. 项目初始化

- [x] 1.1 创建 Go 项目结构（cmd/server, internal, pkg 目录）
- [x] 1.2 初始化 go.mod 并添加依赖（Gin, GORM, go-redis, Zap, Viper）
- [x] 1.3 创建 Makefile（build, run, test, lint, docker 命令）
- [x] 1.4 创建 .gitignore 文件（Go 标准配置）

## 2. 基础设施层

- [x] 2.1 实现配置管理（configs/config.yaml + internal/config/config.go 使用 Viper）
- [x] 2.2 实现 MySQL 连接封装（pkg/database/mysql.go）
- [x] 2.3 实现 Redis 连接封装（pkg/cache/redis.go）
- [x] 2.4 实现数据库模型（internal/model/link.go - ShortLink 结构体）
- [x] 2.5 实现 GORM AutoMigrate 自动建表

## 3. 核心功能 - 短码生成

- [x] 3.1 实现 Snowflake 算法（internal/service/code_generator.go）
- [x] 3.2 实现 Base62 编码转换函数
- [x] 3.3 编写短码生成单元测试（验证唯一性、长度、字符集）

## 4. 核心功能 - URL 缩短

- [x] 4.1 实现 URL 验证工具（internal/util/validator.go）
- [x] 4.2 实现 Repository 层（internal/repository/link_repo.go - Create 方法）
- [x] 4.3 实现 Service 层（internal/service/link_service.go - CreateShortLink 方法）
- [x] 4.4 实现 Handler 层（internal/handler/link_handler.go - CreateLink 方法）
- [x] 4.5 实现统一响应格式（internal/util/response.go）
- [x] 4.6 注册路由（POST /api/v1/links）
- [ ] 4.7 编写创建短链集成测试

## 5. 核心功能 - 短链重定向

- [x] 5.1 实现 Repository 层 GetByCode 方法
- [x] 5.2 实现 Service 层 GetByCode 方法（包含缓存逻辑）
- [x] 5.3 实现缓存策略（Write-through + Cache-Aside + 空值缓存）
- [x] 5.4 实现 Redirect Handler（internal/handler/redirect_handler.go）
- [x] 5.5 注册路由（GET /{code}）
- [ ] 5.6 编写重定向集成测试（缓存命中、缓存未命中、404 场景）

## 6. 中间件

- [x] 6.1 实现请求日志中间件（internal/middleware/logger.go 使用 Zap）
- [x] 6.2 实现异常恢复中间件（internal/middleware/recovery.go）
- [x] 6.3 实现限流中间件（internal/middleware/rate_limit.go - Token Bucket 算法）
- [x] 6.4 注册中间件到 Gin 路由

## 7. 应用入口

- [x] 7.1 实现 main.go（cmd/server/main.go）
- [x] 7.2 初始化所有依赖（Config, DB, Redis, Service, Handler）
- [x] 7.3 配置 Gin 路由和中间件
- [x] 7.4 添加健康检查端点（GET /health）
- [x] 7.5 优雅关闭（监听 SIGTERM/SIGINT）

## 8. 测试

- [x] 8.1 编写 CodeGenerator 单元测试
- [ ] 8.2 编写 URL Validator 单元测试
- [ ] 8.3 编写 LinkService 单元测试（Mock Repository 和 Cache）
- [ ] 8.4 编写 LinkHandler 集成测试（httptest）
- [ ] 8.5 编写 RedirectHandler 集成测试
- [ ] 8.6 运行所有测试并确保通过率 > 90%

## 9. 部署配置

- [x] 9.1 创建 Dockerfile（多阶段构建）
- [x] 9.2 创建 docker-compose.yml（MySQL + Redis + ShortLink）
- [x] 9.3 创建数据库初始化脚本（scripts/init_db.sql）
- [x] 9.4 编写 README.md（项目说明、快速启动、API 文档）

## 10. 性能优化

- [ ] 10.1 配置数据库连接池（MaxOpenConns, MaxIdleConns）
- [ ] 10.2 配置 Redis 连接池（PoolSize）
- [ ] 10.3 性能测试：重定向 P99 < 50ms（使用 wrk/hey）
- [ ] 10.4 性能测试：创建短链 P99 < 100ms
- [ ] 10.5 优化缓存命中率（目标 > 95%）
