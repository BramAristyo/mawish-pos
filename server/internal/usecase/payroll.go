package usecase

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/shopspring/decimal"
)

type PayrollUseCase struct {
	Repo           *repository.PayrollRepository
	AttendanceRepo *repository.AttendanceRepository
	EmployeeRepo   *repository.EmployeeRepository
	AuditLog       *AuditLogUseCase
}

func NewPayrollUseCase(
	repo *repository.PayrollRepository,
	attendanceRepo *repository.AttendanceRepository,
	employeeRepo *repository.EmployeeRepository,
	auditLog *AuditLogUseCase,
) *PayrollUseCase {
	return &PayrollUseCase{
		Repo:           repo,
		AttendanceRepo: attendanceRepo,
		EmployeeRepo:   employeeRepo,
		AuditLog:       auditLog,
	}
}

func (u *PayrollUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.PayrollResponsePagination, error) {
	totalRows, ps, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.PayrollResponsePagination{}, err
	}

	return dto.ToPayrollResponsePagination(ps, req, totalRows), nil
}

func (u *PayrollUseCase) Store(ctx context.Context, req dto.CreatePayrollRequest) (dto.PayrollResponse, error) {
	payrollDomain, err := dto.ToCreatePayrollModel(&req)
	if err != nil {
		return dto.PayrollResponse{}, err
	}

	employee, err := u.EmployeeRepo.FindById(ctx, req.EmployeeID)
	if err != nil {
		return dto.PayrollResponse{}, err
	}

	payrollDomain.BaseSalary = decimal.NewFromFloat(employee.BaseSalary)
	payrollDomain.Employee = &employee

	attendances, err := u.AttendanceRepo.GetByEmployeeID(ctx, req.EmployeeID, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		return dto.PayrollResponse{}, err
	}

	payrollDomain.Calculate(attendances)

	created, err := u.Repo.Store(ctx, payrollDomain)
	if err != nil {
		return dto.PayrollResponse{}, err
	}

	return dto.ToPayrollResponse(created), nil
}
