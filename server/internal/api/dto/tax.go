package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TaxResponse struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Percentage decimal.Decimal `json:"percentage"`
	Status     bool            `json:"status"`
	DeletedAt  *string         `json:"deletedAt,omitempty"`
	CreatedAt  string          `json:"createdAt"`
}

type TaxResponsePagination struct {
	Data []TaxResponse `json:"data"`
	Meta filter.Meta   `json:"meta"`
}

type CreateTaxRequest struct {
	Name       string          `json:"name" binding:"required,min=3,max=100"`
	Percentage decimal.Decimal `json:"percentage" binding:"required,gt=0,lte=100"`
}

type UpdateTaxRequest struct {
	Name       string          `json:"name" binding:"required,min=3,max=100"`
	Percentage decimal.Decimal `json:"percentage" binding:"required,gt=0,lte=100"`
}

func ToCreateTaxModel(req *CreateTaxRequest) domain.Tax {
	return domain.Tax{
		Name:       req.Name,
		Percentage: req.Percentage,
	}
}

func ToUpdateTaxModel(req *UpdateTaxRequest) domain.Tax {
	return domain.Tax{
		Name:       req.Name,
		Percentage: req.Percentage,
	}
}

func ToTaxResponse(t *domain.Tax) TaxResponse {
	resp := TaxResponse{
		ID:         t.ID,
		Name:       t.Name,
		Percentage: t.Percentage,
		Status:     t.Status,
		CreatedAt:  t.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if t.DeletedAt.Valid {
		at := t.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &at
	}

	return resp
}

func ToTaxResponses(ts []domain.Tax) []TaxResponse {
	tsRes := make([]TaxResponse, len(ts))
	for i, t := range ts {
		tsRes[i] = ToTaxResponse(&t)
	}

	return tsRes
}

func ToTaxResponsePagination(t []TaxResponse, p filter.PaginationWithInputFilter, totalRows int64) TaxResponsePagination {
	return TaxResponsePagination{
		Data: t,
		Meta: p.ToMeta(totalRows),
	}
}
