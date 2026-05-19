package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductResponse struct {
	ID             uuid.UUID               `json:"id"`
	CategoryID     *uuid.UUID              `json:"categoryId"`
	Name           string                  `json:"name"`
	CategoryName   *string                 `json:"categoryName"`
	Description    *string                 `json:"description"`
	Price          decimal.Decimal         `json:"price"`
	Cogs           decimal.Decimal         `json:"cogs"`
	ImageURL       *string                 `json:"imageUrl"`
	DeletedAt      *string                 `json:"deletedAt,omitempty"`
	CreatedAt      string                  `json:"createdAt"`
	UpdatedAt      string                  `json:"updatedAt"`
	Category       *CategoryResponse       `json:"category,omitempty"`
	ModifierGroups []ModifierGroupResponse `json:"modifierGroups,omitempty"`
}

type ProductResponsePagination struct {
	Data []ProductResponse `json:"data"`
	Meta filter.Meta       `json:"meta"`
}

type CreateProductRequest struct {
	CategoryID       uuid.UUID       `json:"categoryId" binding:"required,uuid"`
	ModifierGroupIDs []uuid.UUID     `json:"modifierGroupIds" binding:"omitempty,dive,uuid"`
	Name             string          `json:"name" binding:"required,min=3,max=100"`
	Description      *string         `json:"description" binding:"omitempty,max=255"`
	Price            decimal.Decimal `json:"price" binding:"required,gt=0"`
	Cogs             decimal.Decimal `json:"cogs" binding:"required,gt=0"`
	ImageURL         *string         `json:"imageUrl" binding:"omitempty,url"`
}

type UpdateProductRequest struct {
	CategoryID       uuid.UUID       `json:"categoryId" binding:"required,uuid"`
	ModifierGroupIDs []uuid.UUID     `json:"modifierGroupIds" binding:"omitempty,dive,uuid"`
	Name             string          `json:"name" binding:"required,min=3,max=100"`
	Description      *string         `json:"description" binding:"omitempty,max=255"`
	Price            decimal.Decimal `json:"price" binding:"required,gt=0"`
	Cogs             decimal.Decimal `json:"cogs" binding:"required,gt=0"`
	ImageURL         *string         `json:"imageUrl" binding:"omitempty,url"`
}

func ToProductResponse(p *domain.Product) ProductResponse {
	var category *CategoryResponse
	var categoryName *string
	if p.Category != nil {
		c := ToCategoryResponse(p.Category)
		category = &c
		categoryName = &c.Name
	}

	var mgs []ModifierGroupResponse
	if len(p.ProductModifiers) > 0 {
		for _, pm := range p.ProductModifiers {
			mgs = append(mgs, ToModifierGroupResponse(&pm.ModifierGroup))
		}
	}

	resp := ProductResponse{
		ID:             p.ID,
		CategoryID:     &p.CategoryID,
		CategoryName:   categoryName,
		Category:       category,
		ModifierGroups: mgs,
		Name:           p.Name,
		Description:    p.Description,
		Price:          p.Price,
		Cogs:           p.Cogs,
		ImageURL:       p.ImageURL,
		CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if p.DeletedAt.Valid {
		t := p.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &t
	}

	return resp
}

func ToProductResponsePagination(p []ProductResponse, f filter.PaginationWithInputFilter, totalRows int64) ProductResponsePagination {
	return ProductResponsePagination{
		Data: p,
		Meta: f.ToMeta(totalRows),
	}
}

func ToProductModel(req *CreateProductRequest) domain.Product {
	var pms []domain.ProductModifier
	for _, mgID := range req.ModifierGroupIDs {
		pms = append(pms, domain.ProductModifier{
			ModifierGroupID: mgID,
		})
	}

	return domain.Product{
		CategoryID:       req.CategoryID,
		ProductModifiers: pms,
		Name:             req.Name,
		Description:      req.Description,
		Price:            req.Price,
		Cogs:             req.Cogs,
		ImageURL:         req.ImageURL,
	}
}

func ToUpdateProductModel(req *UpdateProductRequest) domain.Product {
	var pms []domain.ProductModifier
	for _, mgID := range req.ModifierGroupIDs {
		pms = append(pms, domain.ProductModifier{
			ModifierGroupID: mgID,
		})
	}

	return domain.Product{
		CategoryID:       req.CategoryID,
		ProductModifiers: pms,
		Name:             req.Name,
		Description:      req.Description,
		Price:            req.Price,
		Cogs:             req.Cogs,
		ImageURL:         req.ImageURL,
	}
}
