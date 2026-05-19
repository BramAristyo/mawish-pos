package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	UseCase *usecase.OrderUseCase
}

func NewOrderHandler(u *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		UseCase: u,
	}
}

func (h *OrderHandler) Paginate(c *gin.Context) {
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

func (h *OrderHandler) FindById(c *gin.Context) {
	id, err := helper.ParseUUID(c)
	if err != nil {
		c.Error(err)
		return
	}

	order, err := h.UseCase.FindById(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, order, "success get order")
}

func (h *OrderHandler) CalculateAll(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.UseCase.CalculateAll(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success calculate order")
}

func (h *OrderHandler) Store(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	created, err := h.UseCase.Store(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, created, "success create order")
}

func (h *OrderHandler) Void(c *gin.Context) {
	var req dto.VoidOrderRequest
	id, err := helper.ParseUUID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.UseCase.Void(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success void order")
}
