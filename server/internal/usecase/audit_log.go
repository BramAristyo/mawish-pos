package usecase

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
)

type AuditLogUseCase struct {
	Repo *repository.AuditLogRepository
}

func NewAuditLogUseCase(repo *repository.AuditLogRepository) *AuditLogUseCase {
	return &AuditLogUseCase{Repo: repo}
}

func (u *AuditLogUseCase) Log(ctx context.Context, log domain.AuditLog) {
	u.Repo.Store(ctx, log)
}
