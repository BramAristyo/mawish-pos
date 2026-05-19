package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type PayrollHandler struct {
	UseCase *usecase.PayrollUseCase
}

func NewPayrollHandler(u *usecase.PayrollUseCase) *PayrollHandler {
	return &PayrollHandler{
		UseCase: u,
	}
}

func (h *PayrollHandler) Paginate(c *gin.Context) {
	var req filter.PaginationWithInputFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	req.DynamicFilter.WithDefaultDateRange()
	req.DynamicFilter.WithDefaultSort()

	res, err := h.UseCase.Paginate(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OKPaginate(c, res.Data, res.Meta)
}

func (h *PayrollHandler) Store(c *gin.Context) {
	var req dto.CreatePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	created, err := h.UseCase.Store(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, created, "success create payroll")
}
