package repository

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	DB *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{DB: db}
}

func (r *AuditLogRepository) Store(ctx context.Context, log domain.AuditLog) {
	r.DB.WithContext(ctx).Create(&log)
}
