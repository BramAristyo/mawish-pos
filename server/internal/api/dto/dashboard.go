package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/shopspring/decimal"
)

type SalesSummaryDashboardResponse struct {
	GrossSales       decimal.Decimal `json:"grossSales"` // penjualan bersih
	NetSales         decimal.Decimal `json:"netSales"`   //
	GrossProfit      decimal.Decimal `json:"grossProfit"`
	TransactionCount int64           `json:"transactionCount"`
	AverageSales     decimal.Decimal `json:"averageSales"`
	GrossMargin      float32         `json:"grossMargin"`
}

func ToSalesSummaryDashboardResponse(s domain.SalesSummaryDashboard) SalesSummaryDashboardResponse {
	return SalesSummaryDashboardResponse{
		GrossSales:       s.GrossSales,
		NetSales:         s.NetSales,
		GrossProfit:      s.GrossProfit,
		TransactionCount: s.TransactionCount,
		AverageSales:     s.AverageSales,
		GrossMargin:      s.GrossMargin,
	}
}
