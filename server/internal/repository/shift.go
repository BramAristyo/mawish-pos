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

type ShiftRepository struct {
	DB *gorm.DB
}

func NewShiftRepository(db *gorm.DB) *ShiftRepository {
	return &ShiftRepository{DB: db}
}

func (r *ShiftRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.Shift, error) {
	var totalRows int64

	allowedFields := map[string]string{
		"created_at": "created_at",
		"opened_at":  "opened_at",
		"closed_at":  "closed_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.Shift{}), req.DynamicFilter, nil, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, []domain.Shift{}, err
	}

	if totalRows == 0 {
		return 0, []domain.Shift{}, nil
	}

	shifts := make([]domain.Shift, 0, req.PaginationInput.PageSize)
	if err := q.
		Preload("OpenedByUser").
		Preload("ClosedByUser").
		Offset(req.Offset()).
		Limit(req.PaginationInput.PageSize).
		Find(&shifts).Error; err != nil {
		return 0, []domain.Shift{}, err
	}

	return totalRows, shifts, nil
}

func (r *ShiftRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Shift, error) {
	var shift domain.Shift
	if err := r.DB.WithContext(ctx).
		Preload("ShiftExpenses").
		Preload("OpenedByUser").
		Preload("ClosedByUser").
		Where("id = ?", id).
		First(&shift).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Shift{}, usecase_errors.NotFound
		}
		return domain.Shift{}, err
	}

	return shift, nil
}

// Opening Shift
func (r *ShiftRepository) Store(ctx context.Context, s *domain.Shift) (domain.Shift, error) {
	if err := r.DB.WithContext(ctx).Create(s).Error; err != nil {
		return domain.Shift{}, err
	}

	return *s, nil
}

func (r *ShiftRepository) Update(ctx context.Context, id uuid.UUID, s *domain.Shift) (domain.Shift, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.Shift
		if err := tx.Where("id = ?", id).Preload("ShiftExpenses").First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return usecase_errors.NotFound
			}
			return err
		}

		updateData := map[string]any{
			"closed_by":    s.ClosedBy,
			"closing_cash": s.ClosingCash,
			"notes":        s.Notes,
			"closed_at":    s.ClosedAt,
			"opened_by":    s.OpenedBy,
			"opening_cash": s.OpeningCash,
		}

		if err := tx.Model(&existing).Updates(updateData).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return domain.Shift{}, err
	}

	return r.FindById(ctx, id)
}

func (r *ShiftRepository) CloseShift(ctx context.Context, id uuid.UUID, s *domain.Shift) (domain.Shift, error) {
	var existing domain.Shift
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Shift{}, usecase_errors.NotFound
		}
		return domain.Shift{}, err
	}

	updateData := map[string]any{
		"closed_by":    s.ClosedBy,
		"closing_cash": s.ClosingCash,
		"notes":        s.Notes,
		"closed_at":    s.ClosedAt,
	}

	if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
		return domain.Shift{}, err
	}

	return existing, nil
}

func (r *ShiftRepository) FindOpenShiftByUserId(ctx context.Context, userId uuid.UUID) (domain.Shift, error) {
	var shift domain.Shift

	if err := r.DB.WithContext(ctx).
		Where("opened_by = ? AND closed_at IS NULL", userId).
		Preload("OpenedByUser").
		Preload("ShiftExpenses").
		First(&shift).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Shift{}, usecase_errors.NoOpenShift
		}
		return domain.Shift{}, err
	}

	return shift, nil
}

// func (r *ShiftRepository) Reconciliation(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.ShiftReconciliaton, error) {
// 	// var totalRows int64

// 	// allowedFields := map[string]string{
// 	// 	"opened_at": "shifts.opened_at",
// 	// }

// 	// q := database.BuildQuery(r.DB.WithContext(ctx).Table("shifts"), req.DynamicFilter, nil, allowedFields)

// 	// q = q.Joins("JOIN users ON shifts.opened_by = users.id")

// 	// if err := q.Count(&totalRows).Error; err != nil {
// 	// 	return 0, nil, err
// 	// }

// 	// if totalRows == 0 {
// 	// 	return 0, nil, nil
// 	// }

// 	// cashPaymentsSub := r.DB.Table("payments").
// 	// 	Select("SUM(payments.amount)").
// 	// 	Joins("JOIN orders ON orders.id = payments.order_id").
// 	// 	Where("orders.shift_id = shifts.id AND payments.method = 'cash' AND orders.status = ?", domain.OrderCompleted)

// 	// shiftRecs := make([]domain.ShiftReconciliaton, 0, req.PaginationInput.PageSize)
// 	// err := q.
// 	// 	Select(`
// 	// 		users.name as cashier_name,
// 	// 		TO_CHAR(shifts.opened_at, 'YYYY-MM-DD HH24:MI:SS') as start_time,
// 	// 		TO_CHAR(shifts.closed_at, 'YYYY-MM-DD HH24:MI:SS') as end_time,
// 	// 		(
// 	// 			shifts.opening_cash
// 	// 			+ COALESCE((?), 0)
// 	// 			+ COALESCE((?), 0)
// 	// 			- COALESCE((?), 0)
// 	// 		) as total_expected,
// 	// 		COALESCE(shifts.closing_cash, 0) as total_actual
// 	// 	`, cashPaymentsSub, cashInSub, cashOutSub).
// 	// 	Offset(req.Offset()).
// 	// 	Limit(req.PaginationInput.PageSize).
// 	// 	Scan(&shiftRecs).Error

// 	// if err != nil {
// 	// 	return 0, nil, err
// 	// }

// 	// return totalRows, shiftRecs, nil
// }
