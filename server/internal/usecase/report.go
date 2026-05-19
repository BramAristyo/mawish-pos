package usecase

import (
	"context"
	"fmt"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
)

type ReportUseCase struct {
	OrderRepo    *repository.OrderRepository
	ShiftRepo    *repository.ShiftRepository
	DiscountRepo *repository.DiscountRepository
}

func NewReportUseCase(
	orderRepo *repository.OrderRepository,
	shiftRepo *repository.ShiftRepository,
	discountRepo *repository.DiscountRepository,
) *ReportUseCase {
	return &ReportUseCase{
		OrderRepo:    orderRepo,
		ShiftRepo:    shiftRepo,
		DiscountRepo: discountRepo,
	}
}

func (u *ReportUseCase) SalesSummary(ctx context.Context, req filter.DynamicFilter) (dto.SalesReportResponse, error) {
	summary, err := u.OrderRepo.SalesSummary(ctx, req)
	if err != nil {
		return dto.SalesReportResponse{}, err
	}

	return dto.ToSalesReportResponse(summary), nil
}

func (u *ReportUseCase) GrossProfit(ctx context.Context, req filter.DynamicFilter) (dto.GrossProfitReportResponse, error) {
	grossProfit, err := u.OrderRepo.GrossProfit(ctx, req)
	if err != nil {
		return dto.GrossProfitReportResponse{}, err
	}

	grossProfit.CalculateGrossProfit()
	grossProfit.PercentageCalculation()

	return dto.ToGrossProfitReportResponse(grossProfit), nil
}

func (u *ReportUseCase) Transactions(ctx context.Context, req filter.PaginationWithInputFilter) (dto.TransactionReportResponsePagination, error) {
	totalRows, ts, err := u.OrderRepo.TransactionReport(ctx, req)
	if err != nil {
		return dto.TransactionReportResponsePagination{}, err
	}

	transactionResponses := make([]dto.TransactionReportResponse, 0, len(ts))
	for _, t := range ts {
		transactionResponses = append(transactionResponses, dto.TransactionReportResponse{
			OrderNumber: t.OrderNumber,
			Time:        t.Time,
			Product:     t.Product,
			Price:       t.Price,
		})
	}

	return dto.TransactionReportResponsePagination{Data: transactionResponses, Meta: req.ToMeta(totalRows)}, nil
}

func (u *ReportUseCase) DiscountUsage(ctx context.Context, req filter.PaginationWithInputFilter) (dto.DiscountReportResponse, error) {
	drs, err := u.DiscountRepo.Usage(ctx, req)
	if err != nil {
		return dto.DiscountReportResponse{}, err
	}

	fmt.Println(drs)

	discountWithFooter := domain.DiscountReportWihFooter{Discounts: drs}
	discountWithFooter.CalculateTotal()

	fmt.Println(discountWithFooter.TotalCount)

	fmt.Println(discountWithFooter.Discounts)

	discountResponses := make([]dto.DiscountReport, len(drs))
	for i, dr := range drs {
		discountResponses[i] = dto.DiscountReport{
			Name:          dr.Name,
			Count:         dr.Count,
			GrossDiscount: dr.GrossDiscount,
		}
	}

	return dto.DiscountReportResponse{
		TotalCount:         discountWithFooter.TotalCount,
		TotalGrossDiscount: discountWithFooter.TotalGrossDiscount,
		Discounts:          discountResponses,
	}, nil

}

// func (u *ReportUseCase) ShiftReconciliation(ctx context.Context, req filter.PaginationWithInputFilter) (dto.ShiftReconciliationtResponsePagination, error) {
// 	totalRows, srs, err := u.ShiftRepo.Reconciliation(ctx, req)
// 	if err != nil {
// 		return dto.ShiftReconciliationtResponsePagination{}, err
// 	}

// 	shiftResponses := make([]dto.ShiftReconciliationResponse, 0, len(srs))
// 	for _, sr := range srs {
// 		sr.CalculateDiff()
// 		shiftResponses = append(shiftResponses, dto.ShiftReconciliationResponse{
// 			CashierName:   sr.CashierName,
// 			StartTime:     sr.StartTime,
// 			EndTime:       sr.EndTime,
// 			TotalExpected: sr.TotalExpected,
// 			TotalActual:   sr.TotalActual,
// 			Difference:    sr.Difference,
// 		})
// 	}

// 	return dto.ShiftReconciliationtResponsePagination{Data: shiftResponses, Meta: req.ToMeta(totalRows)}, nil
// }
