package dto

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PayrollResponse struct {
	ID             uuid.UUID       `json:"id"`
	EmployeeName   string          `json:"employeeName"`
	EmployeeCode   string          `json:"employeeCode"`
	PeriodStart    string          `json:"periodStart"`
	PeriodEnd      string          `json:"periodEnd"`
	BaseSalary     decimal.Decimal `json:"baseSalary"`
	TotalDeduction decimal.Decimal `json:"totalDeduction"`
	NetSalary      decimal.Decimal `json:"netSalary"`
	Notes          *string         `json:"notes"`
	CreatedAt      string          `json:"createdAt"`
}

type PayrollResponsePagination struct {
	Data []PayrollResponse `json:"data"`
	Meta filter.Meta       `json:"meta"`
}

type CreatePayrollRequest struct {
	EmployeeID  uuid.UUID `json:"employeeID" binding:"required,uuid"`
	PeriodStart string    `json:"periodStart" binding:"required"`
	PeriodEnd   string    `json:"periodEnd" binding:"required"`
}

func ToPayrollResponse(p domain.Payroll) PayrollResponse {
	return PayrollResponse{
		ID:             p.ID,
		EmployeeName:   p.Employee.Name,
		EmployeeCode:   p.Employee.Code,
		PeriodStart:    p.PeriodStart.Format("2006-01-02 15:04:05"),
		PeriodEnd:      p.PeriodEnd.Format("2006-01-02 15:04:05"),
		BaseSalary:     p.BaseSalary,
		TotalDeduction: p.TotalDeduction,
		NetSalary:      p.NetSalary,
		Notes:          p.Notes,
		CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToPayrollResponsePagination(ps []domain.Payroll, f filter.PaginationWithInputFilter, totalRows int64) PayrollResponsePagination {
	res := make([]PayrollResponse, len(ps))
	for i, p := range ps {
		res[i] = ToPayrollResponse(p)
	}
	return PayrollResponsePagination{
		Data: res,
		Meta: f.ToMeta(totalRows),
	}
}

func ToCreatePayrollModel(req *CreatePayrollRequest) (domain.Payroll, error) {
	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)

	if err != nil {
		return domain.Payroll{}, err
	}

	periodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)

	if err != nil {
		return domain.Payroll{}, err
	}

	return domain.Payroll{
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		EmployeeID:  req.EmployeeID,
	}, nil
}
