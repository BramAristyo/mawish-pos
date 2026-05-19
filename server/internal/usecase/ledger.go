package usecase

import (
	"context"
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
)

type LedgerUseCase struct {
	Repo *repository.LedgerRepository
}

func NewLedgerUseCase(repo *repository.LedgerRepository) *LedgerUseCase {
	return &LedgerUseCase{
		Repo: repo,
	}
}

func (u *LedgerUseCase) TransactionList(ctx context.Context, req filter.PaginationWithInputFilter) (dto.TransactionSummaryResponse, error) {
	startDate, endDate := u.extractDateRange(req.DynamicFilter)

	summary, err := u.Repo.TransactionSummary(ctx, startDate, endDate)
	if err != nil {
		return dto.TransactionSummaryResponse{}, err
	}

	totalRows, results, err := u.Repo.ReportPaginate(ctx, startDate, endDate, req.PageSize, req.Offset())
	if err != nil {
		return dto.TransactionSummaryResponse{}, err
	}

	transactions := make([]dto.TransactionResponse, 0, len(results))
	for _, r := range results {
		transactions = append(transactions, dto.TransactionResponse{
			TransactionDate: r.TransactionDate,
			ReferenceType:   r.ReferenceType,
			COAName:         r.COA.Name,
			COAType:         r.COA.Type,
			RunningBalance:  r.RunningBalance,
		})
	}

	return dto.TransactionSummaryResponse{
		OpeningBalance: summary.OpeningBalance,
		TotalIncome:    summary.TotalIncome,
		TotalExpense:   summary.TotalExpense,
		Total:          summary.Total,
		Transactions:   transactions,
		Meta:           req.ToMeta(totalRows),
	}, nil
}

func (u *LedgerUseCase) CashFlowStatement(ctx context.Context, req filter.PaginationWithInputFilter) (dto.CashFlowReportResponse, error) {
	startDate, endDate := u.extractDateRange(req.DynamicFilter)

	report, incomes, expenses, err := u.Repo.CashFlowStatement(ctx, startDate, endDate)
	if err != nil {
		return dto.CashFlowReportResponse{}, err
	}

	report.Calculate()

	incomeRes := make([]dto.LedgerResponse, 0, len(incomes))
	for _, i := range incomes {
		incomeRes = append(incomeRes, dto.LedgerResponse{
			ID:              i.ID,
			COAName:         i.COA.Name,
			COAType:         i.COA.Type,
			Amount:          i.Amount,
			ReferenceID:     i.ReferenceID,
			ReferenceType:   i.ReferenceType,
			TransactionDate: i.TransactionDate,
			CreatedAt:       i.CreatedAt,
		})
	}

	expenseRes := make([]dto.LedgerResponse, 0, len(expenses))
	for _, e := range expenses {
		expenseRes = append(expenseRes, dto.LedgerResponse{
			ID:              e.ID,
			COAName:         e.COA.Name,
			COAType:         e.COA.Type,
			Amount:          e.Amount,
			ReferenceID:     e.ReferenceID,
			ReferenceType:   e.ReferenceType,
			TransactionDate: e.TransactionDate,
			CreatedAt:       e.CreatedAt,
		})
	}

	return dto.CashFlowReportResponse{
		OpeningBalance: report.OpeningBalance,
		TotalIncome:    report.TotalIncome,
		TotalExpense:   report.TotalExpense,
		CashFlowAmount: report.CashFlowAmount,
		EndingBalance:  report.Total,
		Incomes:        incomeRes,
		Expenses:       expenseRes,
	}, nil
}

func (u *LedgerUseCase) extractDateRange(df filter.DynamicFilter) (string, string) {
	for _, f := range df.Filter {
		if f.FilterType == filter.DataTypeDate && f.Type == filter.OpInRange {
			if f.From != "" && f.To != "" {
				return f.From, f.To
			}
		}
	}

	// Fallback to current month if not found
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	return startOfMonth.Format("2006-01-02"), now.Format("2006-01-02")
}
