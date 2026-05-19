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

type ShiftScheduleRepository struct {
	DB *gorm.DB
}

func NewShiftScheduleRepository(db *gorm.DB) *ShiftScheduleRepository {
	return &ShiftScheduleRepository{DB: db}
}

func (r *ShiftScheduleRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.ShiftSchedule, error) {
	var ss []domain.ShiftSchedule
	var totalRows int64

	allowedFields := map[string]string{
		"name":       "name",
		"created_at": "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.ShiftSchedule{}), req.DynamicFilter, []string{"name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, []domain.ShiftSchedule{}, err
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&ss).Error; err != nil {
		return 0, []domain.ShiftSchedule{}, err
	}

	return totalRows, ss, nil
}

func (r *ShiftScheduleRepository) FindById(ctx context.Context, id uuid.UUID) (domain.ShiftSchedule, error) {
	var s domain.ShiftSchedule

	if err := r.DB.WithContext(ctx).First(&s, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ShiftSchedule{}, usecase_errors.NotFound
		}
		return domain.ShiftSchedule{}, err
	}

	return s, nil
}

func (r *ShiftScheduleRepository) Store(ctx context.Context, s *domain.ShiftSchedule) (domain.ShiftSchedule, error) {
	if err := r.DB.WithContext(ctx).Create(s).Error; err != nil {
		return domain.ShiftSchedule{}, err
	}

	return *s, nil
}

func (r *ShiftScheduleRepository) Update(ctx context.Context, id uuid.UUID, s *domain.ShiftSchedule) (domain.ShiftSchedule, error) {
	var existing domain.ShiftSchedule
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.ShiftSchedule{}, usecase_errors.NotFound
		}
		return domain.ShiftSchedule{}, err
	}

	updateData := map[string]any{
		"name":                  s.Name,
		"start_time":            s.StartTime,
		"end_time":              s.EndTime,
		"tolerance_minutes":     s.ToleranceMinutes,
		"late_interval_minutes": s.LateIntervalMinutes,
		"late_deduction_amount": s.LateDeductionAmount,
	}

	if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
		return domain.ShiftSchedule{}, err
	}

	return existing, nil
}

func (r *ShiftScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&domain.ShiftSchedule{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *ShiftScheduleRepository) Restore(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).
		Model(&domain.ShiftSchedule{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}
	return result.Error
}

func (r *ShiftScheduleRepository) GetAll(ctx context.Context) ([]domain.ShiftSchedule, error) {
	var ss []domain.ShiftSchedule
	if err := r.DB.WithContext(ctx).Find(&ss).Error; err != nil {
		return nil, err
	}
	return ss, nil
}
