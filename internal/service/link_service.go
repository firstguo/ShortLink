package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/shortlink/shortlink-service/internal/config"
	"github.com/shortlink/shortlink-service/internal/model"
	"github.com/shortlink/shortlink-service/internal/repository"
	"github.com/shortlink/shortlink-service/internal/util"
	"github.com/shortlink/shortlink-service/pkg/cache"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrShortLinkNotFound = errors.New("short link not found")
)

// LinkService 短链服务接口
type LinkService interface {
	CreateShortLink(ctx context.Context, originalURL string) (*model.ShortLinkResponse, error)
	GetByCode(ctx context.Context, code string) (*model.ShortLink, error)
}

type linkService struct {
	repo   repository.LinkRepository
	cache  *redis.Client
	config *config.ShortLinkConfig
}

// NewLinkService 创建短链服务实例
func NewLinkService(repo repository.LinkRepository, cfg *config.ShortLinkConfig) LinkService {
	return &linkService{
		repo:   repo,
		cache:  cache.GetClient(),
		config: cfg,
	}
}

// CreateShortLink 创建短链
func (s *linkService) CreateShortLink(ctx context.Context, originalURL string) (*model.ShortLinkResponse, error) {
	// 1. 验证 URL
	if err := util.ValidateURL(originalURL); err != nil {
		return nil, err
	}

	// 2. 生成短码
	generator, err := NewCodeGenerator(
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		s.config.WorkerID,
		s.config.CodeLength,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create code generator: %w", err)
	}
	code := generator.Generate()

	// 3. 构建短链对象
	link := &model.ShortLink{
		Code:        code,
		OriginalURL: originalURL,
		ShortURL:    fmt.Sprintf("%s/%s", s.config.Domain, code),
		IsEnabled:   true,
	}

	// 4. 保存到数据库
	if err := s.repo.Create(ctx, link); err != nil {
		return nil, fmt.Errorf("failed to save short link: %w", err)
	}

	// 5. 写入缓存（Write-through）
	cacheKey := fmt.Sprintf("shortlink:%s", code)
	cacheData, _ := json.Marshal(map[string]interface{}{
		"id":           link.ID,
		"code":         link.Code,
		"original_url": link.OriginalURL,
		"short_url":    link.ShortURL,
		"is_enabled":   link.IsEnabled,
	})

	// 设置缓存 TTL（24小时 + 随机偏移防雪崩）
	ttl := 24*time.Hour + time.Duration(randInt(0, 8640))*time.Second
	s.cache.Set(ctx, cacheKey, cacheData, ttl)

	// 6. 返回响应
	return &model.ShortLinkResponse{
		Code:        link.Code,
		OriginalURL: link.OriginalURL,
		ShortURL:    link.ShortURL,
		CreatedAt:   link.CreatedAt,
	}, nil
}

// GetByCode 根据短码获取短链（带缓存）
func (s *linkService) GetByCode(ctx context.Context, code string) (*model.ShortLink, error) {
	// 1. 先查缓存
	cacheKey := fmt.Sprintf("shortlink:%s", code)
	cached, err := s.cache.Get(ctx, cacheKey).Result()

	if err == nil {
		// 缓存命中
		if cached == "NULL" {
			return nil, ErrShortLinkNotFound
		}

		// 解析缓存数据
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			link := &model.ShortLink{
				ID:          int64(data["id"].(float64)),
				Code:        data["code"].(string),
				OriginalURL: data["original_url"].(string),
				ShortURL:    data["short_url"].(string),
				IsEnabled:   data["is_enabled"].(bool),
			}
			return link, nil
		}
	}

	// 2. 缓存未命中，查数据库
	link, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 空值缓存（防穿透）
			s.cache.Set(ctx, cacheKey, "NULL", 5*time.Minute)
			return nil, ErrShortLinkNotFound
		}
		return nil, fmt.Errorf("failed to get short link: %w", err)
	}

	// 3. 写入缓存（带随机偏移防雪崩）
	cacheData, _ := json.Marshal(map[string]interface{}{
		"id":           link.ID,
		"code":         link.Code,
		"original_url": link.OriginalURL,
		"short_url":    link.ShortURL,
		"is_enabled":   link.IsEnabled,
	})
	ttl := 24*time.Hour + time.Duration(randInt(0, 8640))*time.Second
	s.cache.Set(ctx, cacheKey, cacheData, ttl)

	return link, nil
}

// randInt 生成 [min, max) 范围内的随机整数
func randInt(min, max int) int {
	return min + int(time.Now().UnixNano()%int64(max-min))
}
