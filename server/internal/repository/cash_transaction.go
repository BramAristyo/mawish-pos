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

type CashTransactionRepository struct {
	DB *gorm.DB
}

func NewCashTransactionRepository(db *gorm.DB) *CashTransactionRepository {
	return &CashTransactionRepository{
		DB: db,
	}
}

func (r *CashTransactionRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.CashTransaction, error) {
	var cts []domain.CashTransaction
	var totalRows int64

	allowedFields := map[string]string{
		"coa_name":   "chart_of_accounts.name",
		"created_at": "cash_transactions.created_at",
	}

	baseQ := r.DB.Model(&domain.CashTransaction{}).
		Joins("JOIN chart_of_accounts ON chart_of_accounts.id = cash_transactions.coa_id").
		Preload("COA")

	q := database.BuildQuery(baseQ, req.DynamicFilter, []string{"chart_of_accounts.name"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, []domain.CashTransaction{}, err
	}

	if err := q.Offset(req.Offset()).Limit(req.PaginationInput.PageSize).Find(&cts).Error; err != nil {
		return 0, []domain.CashTransaction{}, err
	}

	return totalRows, cts, nil
}

func (r *CashTransactionRepository) FindById(ctx context.Context, id uuid.UUID) (domain.CashTransaction, error) {
	var ct domain.CashTransaction

	if err := r.DB.WithContext(ctx).Preload("COA").Preload("Shift").First(&ct, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.CashTransaction{}, usecase_errors.NotFound
		}
		return domain.CashTransaction{}, err
	}

	return ct, nil
}

func (r *CashTransactionRepository) Store(ctx context.Context, ct *domain.CashTransaction) (domain.CashTransaction, error) {
	if err := r.DB.WithContext(ctx).Create(ct).Error; err != nil {
		return domain.CashTransaction{}, err
	}

	return *ct, nil
}

func (r *CashTransactionRepository) Update(ctx context.Context, id uuid.UUID, ct *domain.CashTransaction) (domain.CashTransaction, error) {
	var existing domain.CashTransaction
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.CashTransaction{}, usecase_errors.NotFound
		}
		return domain.CashTransaction{}, err
	}

	updateData := map[string]any{
		"coa_id":      ct.COAID,
		"type":        ct.Type,
		"amount":      ct.Amount,
		"description": ct.Description,
		"date":        ct.Date,
	}

	if err := r.DB.WithContext(ctx).Model(&existing).Updates(updateData).Error; err != nil {
		return domain.CashTransaction{}, err
	}

	return existing, nil
}

func (r *CashTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&domain.CashTransaction{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}

	return result.Error
}
