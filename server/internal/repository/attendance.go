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

type AttendanceRepository struct {
	DB *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{
		DB: db,
	}
}

func (r *AttendanceRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.Attendance, error) {
	var attendances []domain.Attendance
	var totalRows int64

	allowedFields := map[string]string{
		"date":       "date",
		"created_at": "created_at",
	}

	// We might want to allow filtering by employee name/code via join,
	// but for now let's keep it simple or check BuildQuery capabilities.
	// Most repositories use simple fields first.

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.Attendance{}).Preload("Employee").Preload("ShiftSchedule"), req.DynamicFilter, []string{}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Order("date DESC, created_at DESC").Find(&attendances).Error; err != nil {
		return 0, nil, err
	}

	return totalRows, attendances, nil
}

func (r *AttendanceRepository) GetByEmployeeID(ctx context.Context, employeeId uuid.UUID, startPeriod string, endPeriod string) ([]domain.Attendance, error) {
	var attendance []domain.Attendance

	if err := r.DB.WithContext(ctx).
		Where("employee_id = ?", employeeId).
		Where("date BETWEEN ? AND ?", startPeriod, endPeriod).
		Find(&attendance).
		Error; err != nil {
		return []domain.Attendance{}, err
	}

	return attendance, nil
}

func (r *AttendanceRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Attendance, error) {
	var attendance domain.Attendance
	if err := r.DB.WithContext(ctx).
		Preload("Employee").
		Preload("ShiftSchedule").
		First(&attendance, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Attendance{}, usecase_errors.NotFound
		}
		return domain.Attendance{}, err
	}
	return attendance, nil
}

func (r *AttendanceRepository) Store(ctx context.Context, a *domain.Attendance) (domain.Attendance, error) {
	if err := r.DB.WithContext(ctx).Create(a).Error; err != nil {
		return domain.Attendance{}, err
	}
	return *a, nil
}

func (r *AttendanceRepository) Update(ctx context.Context, id uuid.UUID, a *domain.Attendance) (domain.Attendance, error) {
	var existing domain.Attendance
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Attendance{}, usecase_errors.NotFound
		}

		return domain.Attendance{}, err
	}

	updateData := map[string]any { "photo_key": a.PhotoKey }

	if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
		return domain.Attendance{}, err
	}

	return existing, nil
}
