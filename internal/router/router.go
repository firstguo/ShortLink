package router

import (
	"net/http"
	"time"

	"github.com/shortlink/shortlink-service/internal/config"
	"github.com/shortlink/shortlink-service/internal/handler"
	"github.com/shortlink/shortlink-service/internal/middleware"
	"github.com/shortlink/shortlink-service/internal/repository"
	"github.com/shortlink/shortlink-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter 创建并配置路由器
// 接收基础设施层依赖(Repository、Config、Logger),在内部组装完整的依赖链
func NewRouter(
	linkRepo repository.LinkRepository,
	shortLinkConfig *config.ShortLinkConfig,
	rateLimitConfig *config.RateLimitConfig,
	logger *zap.Logger,
) *gin.Engine {
	// 1. 初始化 Service 层
	linkService := service.NewLinkService(linkRepo, shortLinkConfig)

	// 2. 初始化 Handler 层
	linkHandler := handler.NewLinkHandler(linkService)
	redirectHandler := handler.NewRedirectHandler(linkService)

	// 3. 创建 Gin 引擎
	router := gin.New()

	// 4. 注册全局中间件(顺序很重要: Recovery 必须最先)
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.Logger(logger))

	// 可选: 注册限流中间件
	if rateLimitConfig.Enabled {
		router.Use(middleware.RateLimit(rateLimitConfig.Rate, rateLimitConfig.Burst))
	}

	// 5. 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 6. API 路由组
	api := router.Group("/api/v1")
	{
		api.POST("/links", linkHandler.CreateLink)
	}

	// 7. 短链重定向(根路径)
	router.GET("/:code", redirectHandler.Redirect)

	return router
}
