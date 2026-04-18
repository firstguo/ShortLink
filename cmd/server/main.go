package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shortlink/shortlink-service/internal/config"
	"github.com/shortlink/shortlink-service/internal/repository"
	"github.com/shortlink/shortlink-service/internal/router"
	"github.com/shortlink/shortlink-service/pkg/cache"
	"github.com/shortlink/shortlink-service/pkg/database"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 初始化日志
	logger, err := config.InitLogger(&cfg.Log)
	if err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting ShortLink service")

	// 3. 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
	}
	defer database.Close()

	// 4. 自动迁移
	if err := database.AutoMigrate(); err != nil {
		logger.Fatal("Failed to auto migrate", zap.Error(err))
	}
	// 5. 初始化 Redis
	if err := cache.Init(&cfg.Redis); err != nil {
		logger.Fatal("Failed to init Redis", zap.Error(err))
	}
	defer cache.Close()

	// 6. 初始化 Repository
	linkRepo := repository.NewLinkRepository()

	// 7. 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 8. 创建路由(内部会初始化 Service 和 Handler)
	router := router.NewRouter(
		linkRepo,
		&cfg.ShortLink,
		&cfg.RateLimit,
		logger,
	)

	// 9. 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// 10. 启动服务器（优雅关闭）
	go func() {
		logger.Info("Server is running", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
