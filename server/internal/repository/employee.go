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

type EmployeeRepository struct {
	DB *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{
		DB: db,
	}
}

func (r *EmployeeRepository) GetAll(ctx context.Context) ([]domain.Employee, error) {
	var employees []domain.Employee
	if err := r.DB.WithContext(ctx).Order("created_at DESC").Find(&employees).Error; err != nil {
		return []domain.Employee{}, err
	}
	return employees, nil
}

func (r *EmployeeRepository) GetLast(ctx context.Context) (domain.Employee, error) {
	var employee domain.Employee
	if err := r.DB.WithContext(ctx).Unscoped().Order("created_at DESC").First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Employee{}, usecase_errors.NotFound
		}
		return domain.Employee{}, err
	}
	return employee, nil
}

func (r *EmployeeRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.Employee, error) {
	var employees []domain.Employee
	var totalRows int64

	allowedFields := map[string]string{
		"code":        "code",
		"name":        "name",
		"base_salary": "base_salary",
		"created_at":  "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.Employee{}), req.DynamicFilter, []string{"name", "code"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&employees).Error; err != nil {
		return 0, nil, err
	}

	return totalRows, employees, nil
}

func (r *EmployeeRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Employee, error) {
	var employee domain.Employee
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Employee{}, usecase_errors.NotFound
		}
		return domain.Employee{}, err
	}
	return employee, nil
}

func (r *EmployeeRepository) FindByCode(ctx context.Context, code string) (domain.Employee, error) {
	var employee domain.Employee
	if err := r.DB.WithContext(ctx).Where("code = ?", code).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Employee{}, usecase_errors.NotFound
		}
		return domain.Employee{}, err
	}
	return employee, nil
}

func (r *EmployeeRepository) Store(ctx context.Context, employee *domain.Employee) (domain.Employee, error) {
	if err := r.DB.WithContext(ctx).Create(employee).Error; err != nil {
		return domain.Employee{}, err
	}
	return *employee, nil
}

func (r *EmployeeRepository) Update(ctx context.Context, id uuid.UUID, employee *domain.Employee) (domain.Employee, error) {
	var existing domain.Employee
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Employee{}, usecase_errors.NotFound
		}
		return domain.Employee{}, err
	}

	updateData := map[string]interface{}{
		"name":        employee.Name,
		"phone":       employee.Phone,
		"base_salary": employee.BaseSalary,
	}

	if employee.PinHash != "" {
		updateData["pin_hash"] = employee.PinHash
		updateData["has_changed_pin"] = employee.HasChangedPIN
	}

	if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
		return domain.Employee{}, err
	}

	return existing, nil
}

func (r *EmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&domain.Employee{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *EmployeeRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&domain.Employee{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}
