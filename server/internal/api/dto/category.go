package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
)

type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DeletedAt   *string   `json:"deletedAt,omitempty"`
	CreatedAt   string    `json:"createdAt"`
}

type CategoryResponsePagination struct {
	Data []CategoryResponse `json:"data"`
	Meta filter.Meta        `json:"meta"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description"`
}

func ToCreateCategoryModel(req *CreateCategoryRequest) domain.Category {
	return domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}
}

func ToUpdateCategoryModel(req *UpdateCategoryRequest) domain.Category {
	return domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}
}

func ToCategoryResponse(c *domain.Category) CategoryResponse {
	resp := CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if c.DeletedAt.Valid {
		t := c.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &t
	}

	return resp
}

func ToCategoryResponsePagination(c []CategoryResponse, p filter.PaginationWithInputFilter, totalRows int64) CategoryResponsePagination {
	return CategoryResponsePagination{
		Data: c,
		Meta: p.ToMeta(totalRows),
	}
}
