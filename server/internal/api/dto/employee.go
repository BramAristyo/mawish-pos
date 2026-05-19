package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
)

type EmployeeResponse struct {
	ID            uuid.UUID `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Phone         *string   `json:"phone"`
	BaseSalary    float64   `json:"baseSalary"`
	HasChangedPIN bool      `json:"hasChangedPin"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
	DeletedAt     *string   `json:"deletedAt,omitempty"`
}

type EmployeeResponsePagination struct {
	Data []EmployeeResponse `json:"data"`
	Meta filter.Meta        `json:"meta"`
}

type CreateEmployeeRequest struct {
	Name       string   `json:"name" binding:"required,min=2,max=100"`
	Phone      *string  `json:"phone"`
	BaseSalary *float64 `json:"baseSalary"`
	Pin        string   `json:"pin" binding:"required,min=6"`
}

type UpdateEmployeeRequest struct {
	Name       string   `json:"name" binding:"required,min=2,max=100"`
	Phone      *string  `json:"phone"`
	BaseSalary *float64 `json:"baseSalary"`
	Pin        *string  `json:"pin" binding:"omitempty"`
}

func ToEmployeeResponse(e *domain.Employee) EmployeeResponse {
	resp := EmployeeResponse{
		ID:            e.ID,
		Code:          e.Code,
		Name:          e.Name,
		Phone:         e.Phone,
		BaseSalary:    e.BaseSalary,
		HasChangedPIN: e.HasChangedPIN,
		CreatedAt:     e.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     e.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if e.DeletedAt.Valid {
		at := e.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &at
	}

	return resp
}

func ToEmployeeResponsePagination(e []EmployeeResponse, p filter.PaginationWithInputFilter, totalRows int64) EmployeeResponsePagination {
	return EmployeeResponsePagination{
		Data: e,
		Meta: p.ToMeta(totalRows),
	}
}

func ToCreateEmployeeModel(req *CreateEmployeeRequest) domain.Employee {
	return domain.Employee{
		Name:       req.Name,
		Phone:      req.Phone,
		BaseSalary: *req.BaseSalary,
		PinHash:    req.Pin,
	}
}

func ToUpdateEmployeeModel(req *UpdateEmployeeRequest) domain.Employee {
	e := domain.Employee{
		Name:       req.Name,
		Phone:      req.Phone,
		BaseSalary: *req.BaseSalary,
	}
	if req.Pin != nil {
		e.PinHash = *req.Pin
	}
	return e
}
