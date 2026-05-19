package repository

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/persistence/database"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayrollRepository struct {
	DB *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) *PayrollRepository {
	return &PayrollRepository{DB: db}
}

func (r *PayrollRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.Payroll, error) {
	var p []domain.Payroll
	var totalRows int64

	allowedFields := map[string]string{
		"employee_name": "employee.name",
		"created_at":    "payrolls.created_at",
	}

	baseQ := r.DB.Model(&domain.Payroll{}).
		Joins("JOIN employees ON employees.id = payrolls.employee_id").
		Preload("Employee")

	q := database.BuildQuery(baseQ, req.DynamicFilter, []string{"employees.name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, []domain.Payroll{}, err
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&p).Error; err != nil {
		return 0, []domain.Payroll{}, err
	}

	return totalRows, p, nil
}

func (r *PayrollRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Payroll, error) {
	var p domain.Payroll

	if err := r.DB.WithContext(ctx).Preload("Employee").First(&p, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Payroll{}, usecase_errors.NotFound
		}
		return domain.Payroll{}, err
	}

	return p, nil
}

func (r *PayrollRepository) Store(ctx context.Context, p domain.Payroll) (domain.Payroll, error) {
	if err := r.DB.WithContext(ctx).Create(&p).Error; err != nil {
		return domain.Payroll{}, err
	}

	return p, nil
}
