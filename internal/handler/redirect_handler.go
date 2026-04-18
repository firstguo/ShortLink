package handler

import (
	"net/http"

	"github.com/shortlink/shortlink-service/internal/service"

	"github.com/gin-gonic/gin"
)

// RedirectHandler 重定向处理器
type RedirectHandler struct {
	linkService service.LinkService
}

// NewRedirectHandler 创建重定向处理器实例
func NewRedirectHandler(linkService service.LinkService) *RedirectHandler {
	return &RedirectHandler{
		linkService: linkService,
	}
}

// Redirect 处理短链重定向
// GET /{code}
func (h *RedirectHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	if code == "" {
		c.String(http.StatusBadRequest, "Invalid short code")
		return
	}

	// 查询短链
	link, err := h.linkService.GetByCode(c.Request.Context(), code)
	if err != nil {
		if err == service.ErrShortLinkNotFound {
			c.String(http.StatusNotFound, "Short link not found")
			return
		}
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	// 检查是否启用
	if !link.IsEnabled {
		c.String(http.StatusGone, "Short link is disabled")
		return
	}

	// 执行重定向
	c.Redirect(http.StatusFound, link.OriginalURL)
}
