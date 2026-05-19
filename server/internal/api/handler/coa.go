package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type COAHandler struct {
	UseCase *usecase.COAUseCase
}

func NewCOAHandler(u *usecase.COAUseCase) *COAHandler {
	return &COAHandler{
		UseCase: u,
	}
}

func (h *COAHandler) GetAll(c *gin.Context) {
	res, err := h.UseCase.GetAll(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success get all coa")
}

func (h *COAHandler) GetAllOperational(c *gin.Context) {
	res, err := h.UseCase.GetAllOperational(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success get all coa")
}

func (h *COAHandler) Paginate(c *gin.Context) {
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

func (h *COAHandler) FindById(c *gin.Context) {
	id, err := helper.ParseUUID(c)
	if err != nil {
		c.Error(err)
		return
	}
	coa, err := h.UseCase.FindById(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, coa, "success get coa")
}

func (h *COAHandler) Store(c *gin.Context) {
	var req dto.CreateCOARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	created, err := h.UseCase.Store(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, created, "success create coa")
}

func (h *COAHandler) Update(c *gin.Context) {
	var req dto.UpdateCOARequest
	id, err := helper.ParseUUID(c)
	if err != nil {
		c.Error(err)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.UseCase.Update(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success update coa")
}

func (h *COAHandler) Delete(c *gin.Context) {
	id, err := helper.ParseUUID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.UseCase.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	response.OK(c, nil, "success delete coa")
}

func (h *COAHandler) Restore(c *gin.Context) {
	id, err := helper.ParseUUID(c)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := h.UseCase.Restore(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success restore coa")
}
