# API 集成测试

本目录包含 ShortLink 服务的 API 集成测试脚本。

## 测试文件

- `api_integration_test.go` - API 集成测试主文件

## 测试用例

### 1. TestHealthCheck
测试健康检查接口
- **接口**: `GET /health`
- **预期**: 返回 200 OK，状态为 "ok"

### 2. TestCreateShortLink
测试创建短链接口
- **接口**: `POST /api/v1/links`
- **请求体**: `{"original_url": "https://www.example.com/test-page"}`
- **预期**: 返回 200 OK，包含 code、original_url、short_url

### 3. TestCreateShortLinkWithInvalidURL
测试创建短链 - 缺少必要字段
- **接口**: `POST /api/v1/links`
- **请求体**: `{}`
- **预期**: 返回 400 Bad Request

### 4. TestCreateShortLinkWithEmptyURL
测试创建短链 - 空 URL
- **接口**: `POST /api/v1/links`
- **请求体**: `{"original_url": ""}`
- **预期**: 返回 400 Bad Request

### 5. TestRedirectShortLink
测试短链重定向
- **接口**: `GET /{code}`
- **预期**: 返回 302 重定向到原始 URL

### 6. TestRedirectNotFound
测试重定向 - 短链不存在
- **接口**: `GET /notexist123`
- **预期**: 返回 404 Not Found

### 7. TestRedirectInvalidCode
测试重定向 - 无效短码
- **接口**: `GET /`
- **预期**: 验证路由处理

### 8. TestCreateMultipleShortLinks
测试批量创建短链
- **接口**: `POST /api/v1/links`
- **预期**: 成功创建多个短链并能正确查询

## 运行测试

### 运行所有测试
```bash
make test
# 或
go test -v ./test/
```

### 运行单个测试
```bash
go test -v ./test/ -run TestHealthCheck
go test -v ./test/ -run TestCreateShortLink
go test -v ./test/ -run TestRedirect
```

### 运行测试并生成覆盖率报告
```bash
make test
# 会生成 coverage.html 文件
```

## 测试环境要求

测试需要以下服务运行：
- MySQL 8.0 (端口 3306)
- Redis 7 (端口 6380)

可以使用 docker-compose 启动：
```bash
docker-compose up -d
```

## 测试结果

所有测试通过后，您应该看到：
```
=== RUN   TestHealthCheck
--- PASS: TestHealthCheck (0.00s)
=== RUN   TestCreateShortLink
--- PASS: TestCreateShortLink (0.02s)
=== RUN   TestCreateShortLinkWithInvalidURL
--- PASS: TestCreateShortLinkWithInvalidURL (0.00s)
=== RUN   TestCreateShortLinkWithEmptyURL
--- PASS: TestCreateShortLinkWithEmptyURL (0.00s)
=== RUN   TestRedirectShortLink
--- PASS: TestRedirectShortLink (0.01s)
=== RUN   TestRedirectNotFound
--- PASS: TestRedirectNotFound (0.00s)
=== RUN   TestRedirectInvalidCode
--- PASS: TestRedirectInvalidCode (0.00s)
=== RUN   TestCreateMultipleShortLinks
--- PASS: TestCreateMultipleShortLinks (0.04s)
PASS
```

## 注意事项

1. 测试会连接真实的 MySQL 和 Redis 服务
2. 测试会在数据库中创建真实的记录
3. 测试使用 `gin.TestMode` 模式运行
4. 确保在运行测试前，MySQL 和 Redis 服务已经启动
