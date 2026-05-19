package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/shopspring/decimal"
)

type SalesReportResponse struct {
	GrossSales decimal.Decimal `json:"grossSales"`
	Discounts  decimal.Decimal `json:"discounts"`
	NetSales   decimal.Decimal `json:"netSales"`
	Gratuity   decimal.Decimal `json:"gratuity"`
	Tax        decimal.Decimal `json:"tax"`
	Total      decimal.Decimal `json:"total"`
}

func ToSalesReportResponse(s domain.SalesSummary) SalesReportResponse {
	return SalesReportResponse{
		GrossSales: s.GrossSales,
		Discounts:  s.Discounts,
		NetSales:   s.NetSales,
		Gratuity:   s.Gratuity,
		Tax:        s.Tax,
		Total:      s.Total,
	}
}

type GrossProfitReportResponse struct {
	GrossSales            decimal.Decimal `json:"grossSales"`
	Discounts             decimal.Decimal `json:"discounts"`
	NetSales              decimal.Decimal `json:"netSales"`
	NetSalesPercentage    decimal.Decimal `json:"netSalesPercentage"`
	Cogs                  decimal.Decimal `json:"cogs"`
	CogsPercetage         decimal.Decimal `json:"cogsPercentage"`
	GrossProfit           decimal.Decimal `json:"grossProfit"`
	GrossProfitPercentage decimal.Decimal `json:"grossProfitPercentage"`
}

func ToGrossProfitReportResponse(gp domain.GrossProfit) GrossProfitReportResponse {
	return GrossProfitReportResponse{
		GrossSales:            gp.GrossSales,
		Discounts:             gp.Discounts,
		NetSales:              gp.NetSales,
		NetSalesPercentage:    gp.NetSalesPercentage,
		Cogs:                  gp.Cogs,
		CogsPercetage:         gp.CogsPercentage,
		GrossProfit:           gp.GrossProfit,
		GrossProfitPercentage: gp.GrossProfitPercentage,
	}
}

type DiscountReport struct {
	Name          string          `json:"name"`
	Count         int64           `json:"count"`
	GrossDiscount decimal.Decimal `json:"grossDiscount"`
	// Discount      decimal.Decimal `json:"discount"`
}
type DiscountReportResponse struct {
	TotalCount         int64            `json:"totalCount"`
	TotalGrossDiscount decimal.Decimal  `json:"totalGrossDiscount"`
	Discounts          []DiscountReport `json:"discounts"`
}

type DiscountReportResponsePagination struct {
	Data []DiscountReportResponse `json:"data"`
	Meta filter.Meta              `json:"meta"`
}

type TransactionReportResponse struct {
	OrderNumber string          `json:"orderNumber"`
	Time        string          `json:"time"`
	Product     string          `json:"product"`
	Price       decimal.Decimal `json:"price"`
}
type TransactionReportResponsePagination struct {
	Data []TransactionReportResponse `json:"data"`
	Meta filter.Meta                 `json:"meta"`
}

type ShiftReconciliationResponse struct {
	CashierName   string          `json:"cashierName"`
	StartTime     string          `json:"startTime"`
	EndTime       *string         `json:"endTime,omitempty"`
	TotalExpected decimal.Decimal `json:"totalExpected"`
	TotalActual   decimal.Decimal `json:"totalActual"`
	Difference    decimal.Decimal `json:"difference"`
}
type ShiftReconciliationtResponsePagination struct {
	Data []ShiftReconciliationResponse `json:"data"`
	Meta filter.Meta                   `json:"meta"`
}
