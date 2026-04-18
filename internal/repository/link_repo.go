package repository

import (
	"context"

	"github.com/shortlink/shortlink-service/internal/model"
	"github.com/shortlink/shortlink-service/pkg/database"

	"gorm.io/gorm"
)

// LinkRepository 短链仓储接口
type LinkRepository interface {
	Create(ctx context.Context, link *model.ShortLink) error
	GetByCode(ctx context.Context, code string) (*model.ShortLink, error)
}

type linkRepository struct {
	db *gorm.DB
}

// NewLinkRepository 创建短链仓储实例
func NewLinkRepository() LinkRepository {
	return &linkRepository{
		db: database.DB,
	}
}

// Create 创建短链记录
func (r *linkRepository) Create(ctx context.Context, link *model.ShortLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

// GetByCode 根据短码查询短链
func (r *linkRepository) GetByCode(ctx context.Context, code string) (*model.ShortLink, error) {
	var link model.ShortLink
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}
