package repository

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type LedgerRepository struct {
	DB *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) *LedgerRepository {
	return &LedgerRepository{
		DB: db,
	}
}

func (r *LedgerRepository) ReportPaginate(ctx context.Context, startDate string, endDate string, limit int, offset int) (int64, []domain.LedgerWithBalance, error) {
	var totalRows int64
	var results []domain.LedgerWithBalance

	if err := r.DB.WithContext(ctx).Model(&domain.Ledger{}).
		Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
		Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	err := r.DB.WithContext(ctx).
		Model(&domain.Ledger{}).
		Select("ledgers.*, SUM(CASE WHEN ca.type = 'in' THEN ledgers.amount WHEN ca.type = 'out' THEN -ledgers.amount ELSE 0 END) OVER (ORDER BY ledgers.transaction_date ASC, ledgers.created_at ASC) AS running_balance").
		Joins("JOIN chart_of_accounts ca ON ca.id = ledgers.coa_id").
		Where("ledgers.transaction_date BETWEEN ? AND ?", startDate, endDate).
		Preload("COA").
		Order("transaction_date ASC, created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&results).Error

	if err != nil {
		return 0, nil, err
	}

	return totalRows, results, nil
}

func (r *LedgerRepository) TransactionSummary(ctx context.Context, startDate string, endDate string) (domain.TransactionSummary, error) {
	q := `
		SELECT
			COALESCE(SUM(
				CASE
					WHEN ca.type = 'in' THEN l.amount
					WHEN ca.type = 'out' THEN -l.amount
				END
			) FILTER (WHERE l.transaction_date < ?), 0) AS opening_balance,

			COALESCE(SUM(l.amount) FILTER (where ca.type = 'in' AND l.transaction_date BETWEEN ? AND ?), 0) AS total_income,
			COALESCE(SUM(l.amount) FILTER (where ca.type = 'out' AND l.transaction_date BETWEEN ? AND ?), 0) AS total_expense,

			COALESCE(SUM(
				CASE
					WHEN ca.type = 'in' THEN l.amount
					WHEN ca.type = 'out' THEN -l.amount
				END
			) FILTER (WHERE l.transaction_date BETWEEN ? AND ?), 0) as total

			FROM ledgers l
			JOIN chart_of_accounts ca ON ca.id = l.coa_id
	`
	var summary domain.TransactionSummary
	err := r.DB.WithContext(ctx).Raw(q,
		startDate,
		startDate, endDate,
		startDate, endDate,
		startDate, endDate,
	).Scan(&summary).Error

	if err != nil {
		return domain.TransactionSummary{}, err
	}

	return summary, nil

}

func (r *LedgerRepository) CashFlowStatement(ctx context.Context, startDate string, endDate string) (domain.CashFlowReport, []domain.Ledger, []domain.Ledger, error) {

	var summary domain.CashFlowReport
	var incomes []domain.Ledger
	var expenses []domain.Ledger

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		summaryQ := `
			SELECT
				COALESCE(SUM(
					CASE
						WHEN ca.type = 'in' THEN l.amount
						WHEN ca.type = 'out' THEN -l.amount
						ELSE 0
					END
				) FILTER (WHERE l.transaction_date < ?), 0) AS opening_balance,
				COALESCE(SUM(l.amount) FILTER (WHERE ca.type = 'in' AND l.transaction_date BETWEEN ? AND ?), 0) AS total_income,
				COALESCE(SUM(l.amount) FILTER (WHERE ca.type = 'out' AND l.transaction_date BETWEEN ? AND ?), 0) AS total_expense
			FROM ledgers l
			JOIN chart_of_accounts ca ON ca.id = l.coa_id
		`

		return r.DB.WithContext(gctx).Raw(summaryQ,
			startDate,
			startDate, endDate,
			startDate, endDate,
		).Scan(&summary).Error
	})

	g.Go(func() error {
		if err := r.DB.WithContext(gctx).
			Model(&domain.Ledger{}).
			Joins("JOIN chart_of_accounts ca ON ca.id = ledgers.coa_id").
			Where("ledgers.transaction_date BETWEEN ? AND ?", startDate, endDate).
			Where("ca.type = ?", "in").
			Preload("COA").
			Order("transaction_date ASC, created_at ASC").
			Find(&incomes).Error; err != nil {
			return err
		}

		return r.DB.WithContext(gctx).
			Model(&domain.Ledger{}).
			Joins("JOIN chart_of_accounts ca ON ca.id = ledgers.coa_id").
			Where("ledgers.transaction_date BETWEEN ? AND ?", startDate, endDate).
			Where("ca.type = ?", "out").
			Preload("COA").
			Order("transaction_date ASC, created_at ASC").
			Find(&expenses).Error
	})

	if err := g.Wait(); err != nil {
		return domain.CashFlowReport{}, nil, nil, err
	}

	return summary, incomes, expenses, nil
}

func (r *LedgerRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Ledger, error) {
	var l domain.Ledger

	if err := r.DB.WithContext(ctx).First(&l, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Ledger{}, usecase_errors.NotFound
		}
		return domain.Ledger{}, err
	}

	return l, nil
}

func (r *LedgerRepository) Store(ctx context.Context, ledger domain.Ledger) (domain.Ledger, error) {
	if err := r.DB.WithContext(ctx).Create(&ledger).Error; err != nil {
		return domain.Ledger{}, err
	}

	return ledger, nil
}

func (r *LedgerRepository) Delete(ctx context.Context, refId uuid.UUID, refType domain.ReferenceType) error {
	result := r.DB.WithContext(ctx).Delete(&domain.Ledger{}, "reference_id = ? AND reference_type = ?", refId, refType)
	if result.RowsAffected == 0 {
		return usecase_errors.NotFound
	}

	return result.Error
}
