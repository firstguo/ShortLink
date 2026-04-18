package database

import (
	"fmt"

	"github.com/shortlink/shortlink-service/internal/model"
)

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if err := DB.AutoMigrate(
		&model.ShortLink{},
	); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	fmt.Println("Database migration completed")
	return nil
}
