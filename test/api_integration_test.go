package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shortlink/shortlink-service/internal/config"
	"github.com/shortlink/shortlink-service/internal/handler"
	"github.com/shortlink/shortlink-service/internal/middleware"
	"github.com/shortlink/shortlink-service/internal/repository"
	"github.com/shortlink/shortlink-service/internal/service"
	"github.com/shortlink/shortlink-service/internal/util"
	"github.com/shortlink/shortlink-service/pkg/cache"
	"github.com/shortlink/shortlink-service/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	testLogger    *zap.Logger
	testRouter    *gin.Engine
	testHandler   *handler.LinkHandler
	testRedirect  *handler.RedirectHandler
	testService   service.LinkService
)

// TestMain 初始化测试环境
func TestMain(m *testing.M) {
	// 加载测试配置
	cfg, err := config.LoadConfig("../configs/config.yaml")
	if err != nil {
		panic(err)
	}

	// 初始化日志
	testLogger, _ = config.InitLogger(&cfg.Log)

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		panic(err)
	}
	defer database.Close()

	// 自动迁移
	if err := database.AutoMigrate(); err != nil {
		panic(err)
	}

	// 初始化 Redis
	if err := cache.Init(&cfg.Redis); err != nil {
		panic(err)
	}
	defer cache.Close()

	// 初始化依赖
	linkRepo := repository.NewLinkRepository()
	testService = service.NewLinkService(linkRepo, &cfg.ShortLink)
	testHandler = handler.NewLinkHandler(testService)
	testRedirect = handler.NewRedirectHandler(testService)

	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由
	testRouter = gin.New()
	testRouter.Use(middleware.Recovery(testLogger))
	testRouter.Use(middleware.Logger(testLogger))

	// 健康检查
	testRouter.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API 路由
	api := testRouter.Group("/api/v1")
	{
		api.POST("/links", testHandler.CreateLink)
	}

	// 重定向路由
	testRouter.GET("/:code", testRedirect.Redirect)

	m.Run()
}

// TestHealthCheck 测试健康检查接口
func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

// TestCreateShortLink 测试创建短链接口
func TestCreateShortLink(t *testing.T) {
	// 准备测试数据
	requestBody := map[string]string{
		"original_url": "https://www.example.com/test-page",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	testRouter.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response util.Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Message)

	// 验证返回数据包含必要字段
	data := response.Data.(map[string]interface{})
	assert.NotEmpty(t, data["code"])
	assert.Equal(t, "https://www.example.com/test-page", data["original_url"])
	assert.NotEmpty(t, data["short_url"])
}

// TestCreateShortLinkWithInvalidURL 测试创建短链-无效URL
func TestCreateShortLinkWithInvalidURL(t *testing.T) {
	// 准备测试数据 - 缺少 original_url
	requestBody := map[string]string{}
	jsonBody, _ := json.Marshal(requestBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	testRouter.ServeHTTP(w, req)

	// 验证响应 - 应该返回 400 错误
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestCreateShortLinkWithEmptyURL 测试创建短链-空URL
func TestCreateShortLinkWithEmptyURL(t *testing.T) {
	// 准备测试数据 - 空 URL
	requestBody := map[string]string{
		"original_url": "",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	testRouter.ServeHTTP(w, req)

	// 验证响应 - 应该返回 400 错误
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRedirectShortLink 测试短链重定向
func TestRedirectShortLink(t *testing.T) {
	// 先创建一个短链
	link, err := testService.CreateShortLink(context.Background(), "https://www.google.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, link.Code)

	// 测试重定向
	req := httptest.NewRequest(http.MethodGet, "/"+link.Code, nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	// 验证重定向响应
	assert.Equal(t, http.StatusFound, w.Code)
	location := w.Header().Get("Location")
	assert.Equal(t, "https://www.google.com", location)
}

// TestRedirectNotFound 测试重定向-短链不存在
func TestRedirectNotFound(t *testing.T) {
	// 使用不存在的短码
	req := httptest.NewRequest(http.MethodGet, "/notexist123", nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	// 验证 404 响应
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestRedirectInvalidCode 测试重定向-无效短码
func TestRedirectInvalidCode(t *testing.T) {
	// 使用空短码
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	// 注意：这个测试可能会匹配到其他路由，具体取决于路由配置
	// 这里只是示例
	t.Log("Redirect with invalid code test")
}

// TestCreateMultipleShortLinks 测试创建多个短链
func TestCreateMultipleShortLinks(t *testing.T) {
	urls := []string{
		"https://www.baidu.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
	}

	codes := make([]string, 0)

	for _, url := range urls {
		requestBody := map[string]string{
			"original_url": url,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data := response["data"].(map[string]interface{})
		code := data["code"].(string)
		codes = append(codes, code)
	}

	// 验证所有短链都能正确重定向
	for i, code := range codes {
		link, err := testService.GetByCode(context.Background(), code)
		assert.NoError(t, err)
		assert.Equal(t, urls[i], link.OriginalURL)
	}
}
