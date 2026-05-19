package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type LedgerHandler struct {
	usecase *usecase.LedgerUseCase
}

func NewLedgerHandler(u *usecase.LedgerUseCase) *LedgerHandler {
	return &LedgerHandler{usecase: u}
}

func (h *LedgerHandler) TransactionList(c *gin.Context) {
	var req filter.PaginationWithInputFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	req.DynamicFilter.WithDefaultSort()
	req.DynamicFilter.WithDefaultDateRange()

	if err := req.DynamicFilter.ValidateHasDateRange(); err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.TransactionList(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "Transaction list retrieved successfully")
}

func (h *LedgerHandler) CashFlowStatement(c *gin.Context) {
	var req filter.PaginationWithInputFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	req.DynamicFilter.WithDefaultSort()
	req.DynamicFilter.WithDefaultDateRange()

	if err := req.DynamicFilter.ValidateHasDateRange(); err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.CashFlowStatement(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "Cash flow statement retrieved successfully")
}
