package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type ShiftHandler struct {
	UseCase *usecase.ShiftUseCase
}

func NewShiftHandler(u *usecase.ShiftUseCase) *ShiftHandler {
	return &ShiftHandler{
		UseCase: u,
	}
}

func (h *ShiftHandler) Paginate(c *gin.Context) {
	var req filter.PaginationWithInputFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	req.DynamicFilter.WithDefaultSort()

	res, err := h.UseCase.Paginate(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OKPaginate(c, res.Data, res.Meta)
}

func (h *ShiftHandler) FindById(c *gin.Context) {
	id, err := helper.ParseUUID(c)
	if err != nil {
		c.Error(err)
		return
	}
	shift, err := h.UseCase.FindById(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, shift, "success get shift")
}

func (h *ShiftHandler) OpenShift(c *gin.Context) {
	var req dto.OpenShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.UseCase.OpenShift(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, res, "shift opened successfully")
}

func (h *ShiftHandler) CloseShift(c *gin.Context) {
	var req dto.CloseShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.UseCase.CloseShift(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "shift closed successfully")
}

func (h *ShiftHandler) FindOpenShiftByCurrentUser(c *gin.Context) {

	res, err := h.UseCase.FindOpenShiftByCurrent(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "successfully get active shift")
}

func (h *ShiftHandler) UpsertExpenses(c *gin.Context) {
	var req dto.UpsertShiftExpensesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.UseCase.UpsertExpenses(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "expenses updated successfully")
}
