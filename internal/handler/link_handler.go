package handler

import (
	"github.com/shortlink/shortlink-service/internal/service"
	"github.com/shortlink/shortlink-service/internal/util"

	"github.com/gin-gonic/gin"
)

// LinkHandler 短链处理器
type LinkHandler struct {
	linkService service.LinkService
}

// NewLinkHandler 创建短链处理器实例
func NewLinkHandler(linkService service.LinkService) *LinkHandler {
	return &LinkHandler{
		linkService: linkService,
	}
}

// CreateLink 创建短链
// POST /api/v1/links
func (h *LinkHandler) CreateLink(c *gin.Context) {
	var req struct {
		OriginalURL string `json:"original_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// 调用服务层创建短链
	result, err := h.linkService.CreateShortLink(c.Request.Context(), req.OriginalURL)
	if err != nil {
		// 判断错误类型
		if validationErr, ok := err.(*util.ValidationError); ok {
			util.BadRequest(c, validationErr.Message)
			return
		}
		util.InternalServerError(c, "Failed to create short link")
		return
	}

	util.Success(c, result)
}
