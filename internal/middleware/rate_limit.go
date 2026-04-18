package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket 令牌桶限流器
type TokenBucket struct {
	rate       int           // 每秒产生的令牌数
	burst      int           // 桶容量
	tokens     float64       // 当前令牌数
	lastTime   time.Time     // 上次更新时间
	mu         sync.Mutex
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(rate, burst int) *TokenBucket {
	return &TokenBucket{
		rate:     rate,
		burst:    burst,
		tokens:   float64(burst),
		lastTime: time.Now(),
	}
}

// Allow 检查是否允许请求
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastTime).Seconds()
	tb.lastTime = now

	// 添加新令牌
	tb.tokens += elapsed * float64(tb.rate)
	if tb.tokens > float64(tb.burst) {
		tb.tokens = float64(tb.burst)
	}

	// 检查是否有足够的令牌
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}

// RateLimit 限流中间件
func RateLimit(rate, burst int) gin.HandlerFunc {
	limiter := NewTokenBucket(rate, burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
