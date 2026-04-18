## Why

ShortLink 项目目前只有需求文档（REQUIREMENTS.md）和技术设计文档（TECH_DESIGN.md），但没有任何可执行的代码。需要从零开始实现 MVP 版本的核心功能：URL 缩短和短链重定向，以验证技术方案并交付可用服务。

## What Changes

- 实现完整的 Go 项目结构（cmd、internal、pkg 分层架构）
- 实现 URL 缩短功能（POST /api/v1/links）
- 实现短链重定向功能（GET /{code}，302 重定向）
- 集成 MySQL 数据库存储短链映射
- 集成 Redis 缓存提升重定向性能（P99 < 50ms）
- 实现 Snowflake + 62 进制短码生成算法
- 实现中间件：请求日志、异常恢复、限流（Token Bucket）
- 提供 Docker 部署配置和 Makefile
- 编写单元测试和集成测试

## Capabilities

### New Capabilities
- `shortlink-create`: URL 缩短功能，包括 URL 验证、短码生成、数据库存储、缓存写入
- `shortlink-redirect`: 短链重定向功能，包括缓存查询、数据库降级、302 重定向、缓存穿透防护
- `code-generation`: Snowflake 算法 + 62 进制编码生成唯一短码
- `cache-management`: Redis 缓存策略，包括 Write-through、空值缓存、随机 TTL

### Modified Capabilities
<!-- 无现有规格需要修改 -->

## Impact

- **新增代码**: 完整的 Go 项目（约 2000 行代码）
- **新增依赖**: Gin, GORM, go-redis, Zap, Viper
- **基础设施**: 需要 MySQL 8.0+ 和 Redis 7.0+
- **API**: 新增 2 个 RESTful 端点
- **部署**: Docker 容器化部署
