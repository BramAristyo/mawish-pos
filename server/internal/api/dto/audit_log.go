package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
)

type AuditLogResponse struct {
	ID          uuid.UUID          `json:"id"`
	UserID      uuid.UUID          `json:"userId"`
	Action      domain.AuditAction `json:"action"`
	Entity      domain.AuditEntity `json:"entity"`
	EntityId    any                `json:"entityId"`
	Description string             `json:"description"`
	CreatedAt   string             `json:"createdAt"`
	User        UserResponse       `json:"user"`
}

type DailyAuditLogReponse struct {
	Date      string             `json:"date"`
	AuditLogs []AuditLogResponse `json:"auditLogs"`
}

func ToAuditLogResponse(a *domain.AuditLog) AuditLogResponse {
	return AuditLogResponse{
		ID:          a.ID,
		UserID:      a.UserID,
		Action:      a.Action,
		Entity:      a.Entity,
		EntityId:    a.EntityID,
		Description: a.Description,
		CreatedAt:   a.CreatedAt.Format("2006-01-02 15:04:05"),
		User:        ToUserResponse(a.User),
	}
}
