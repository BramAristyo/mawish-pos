package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	usecase *usecase.ReportUseCase
}

func NewReportHandler(u *usecase.ReportUseCase) *ReportHandler {
	return &ReportHandler{usecase: u}
}

func (h *ReportHandler) SalesSummary(c *gin.Context) {
	var req filter.DynamicFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	req.WithDefaultDateRange()

	res, err := h.usecase.SalesSummary(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "Sales summary retrieved successfully")
}

func (h *ReportHandler) GrossProfit(c *gin.Context) {
	var req filter.DynamicFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	req.WithDefaultDateRange()

	res, err := h.usecase.GrossProfit(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "Gross profit retrieved successfully")
}

func (h *ReportHandler) Transactions(c *gin.Context) {
	var req filter.PaginationWithInputFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	req.DynamicFilter.WithDefaultSort()
	req.DynamicFilter.WithDefaultDateRange()

	res, err := h.usecase.Transactions(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "Transactions retrieved successfully")
}

func (h *ReportHandler) DiscountUsage(c *gin.Context) {
	var req filter.PaginationWithInputFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	req.DynamicFilter.WithDefaultSort()
	req.DynamicFilter.WithDefaultDateRange()

	res, err := h.usecase.DiscountUsage(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "successfully get discounts Report")
}

// func (h *ReportHandler) ShiftReconciliation(c *gin.Context) {
// 	var req filter.PaginationWithInputFilter
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	req.DynamicFilter.WithDefaultSort()
// 	req.DynamicFilter.WithDefaultDateRange()

// 	res, err := h.usecase.ShiftReconciliation(c.Request.Context(), req)
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	response.OKPaginate(c, res.Data, res.Meta)
// }
