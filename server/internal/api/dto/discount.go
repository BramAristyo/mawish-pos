package dto

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DiscountResponse struct {
	ID        uuid.UUID             `json:"id"`
	Name      string                `json:"name"`
	Type      domain.AdjustmentType `json:"type"`
	Value     decimal.Decimal       `json:"value"`
	StartDate *string               `json:"startDate"`
	EndDate   *string               `json:"endDate"`
	DeletedAt *string               `json:"deletedAt,omitempty"`
	CreatedAt string                `json:"createdAt"`
}

type DiscountResponsePagination struct {
	Data []DiscountResponse `json:"data"`
	Meta filter.Meta        `json:"meta"`
}

type CreateDiscountRequest struct {
	Name      string                `json:"name" binding:"required,min=3,max=100"`
	Type      domain.AdjustmentType `json:"type" binding:"required,oneof=percentage fixed"`
	Value     decimal.Decimal       `json:"value" binding:"required,gt=0"`
	StartDate *string               `json:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate   *string               `json:"endDate" binding:"omitempty,datetime=2006-01-02"`
}

type UpdateDiscountRequest struct {
	Name      string                `json:"name" binding:"required,min=3,max=100"`
	Type      domain.AdjustmentType `json:"type" binding:"required,oneof=percentage fixed"`
	Value     decimal.Decimal       `json:"value" binding:"required,gt=0"`
	StartDate *string               `json:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate   *string               `json:"endDate" binding:"omitempty,datetime=2006-01-02"`
}

func ToCreateDiscountModel(req *CreateDiscountRequest) domain.Discount {
	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *req.StartDate)
		startDate = &t
	}
	if req.EndDate != nil {
		t, _ := time.Parse("2006-01-02", *req.EndDate)
		endDate = &t
	}

	return domain.Discount{
		Name:      req.Name,
		Type:      req.Type,
		Value:     req.Value,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

func ToUpdateDiscountModel(req *UpdateDiscountRequest) domain.Discount {
	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t, _ := time.Parse("2006-01-02", *req.StartDate)
		startDate = &t
	}
	if req.EndDate != nil {
		t, _ := time.Parse("2006-01-02", *req.EndDate)
		endDate = &t
	}

	return domain.Discount{
		Name:      req.Name,
		Type:      req.Type,
		Value:     req.Value,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

func ToDiscountResponse(d *domain.Discount) DiscountResponse {
	var startDate, endDate *string
	if d.StartDate != nil {
		s := d.StartDate.Format("2006-01-02")
		startDate = &s
	}
	if d.EndDate != nil {
		s := d.EndDate.Format("2006-01-02")
		endDate = &s
	}

	resp := DiscountResponse{
		ID:        d.ID,
		Name:      d.Name,
		Type:      d.Type,
		Value:     d.Value,
		StartDate: startDate,
		EndDate:   endDate,
		CreatedAt: d.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if d.DeletedAt.Valid {
		at := d.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &at
	}

	return resp
}

func ToDiscountResponses(ds []domain.Discount) []DiscountResponse {
	dsRes := make([]DiscountResponse, len(ds))
	for i, d := range ds {
		dsRes[i] = ToDiscountResponse(&d)
	}

	return dsRes
}

func ToDiscountResponsePagination(d []DiscountResponse, p filter.PaginationWithInputFilter, totalRows int64) DiscountResponsePagination {
	return DiscountResponsePagination{
		Data: d,
		Meta: p.ToMeta(totalRows),
	}
}
