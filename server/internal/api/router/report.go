package router

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func ReportRoutes(r *gin.RouterGroup, h *handler.ReportHandler) {
	r.GET("/sales-summary", h.SalesSummary)
	r.GET("/gross-profit", h.GrossProfit)
	r.GET("/transactions", h.Transactions)
	r.GET("/discount-usage", h.DiscountUsage)
	// r.GET("/shift-reconciliation", h.ShiftReconciliation)
}
