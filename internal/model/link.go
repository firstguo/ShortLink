package model

import "time"

// ShortLink 短链模型
type ShortLink struct {
	ID          int64     `gorm:"primaryKey;autoIncrement;comment:'主键ID'" json:"id"`
	Code        string    `gorm:"type:varchar(10);uniqueIndex;not null;comment:'短码'" json:"code"`
	OriginalURL string    `gorm:"type:varchar(2048);not null;comment:'原始URL'" json:"original_url"`
	ShortURL    string    `gorm:"type:varchar(256);comment:'短链URL'" json:"short_url"`
	IsEnabled   bool      `gorm:"default:true;comment:'是否启用'" json:"is_enabled"`
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:'创建时间'" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:'更新时间'" json:"updated_at"`
}

// TableName 指定表名
func (ShortLink) TableName() string {
	return "short_links"
}

// CreateShortLinkRequest 创建短链请求
type CreateShortLinkRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
}

// ShortLinkResponse 短链响应
type ShortLinkResponse struct {
	Code        string    `json:"code"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}
