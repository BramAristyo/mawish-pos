package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	UseCase *usecase.CategoryUseCase
}

func NewCategoryHandler(u *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		UseCase: u,
	}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	res, err := h.UseCase.GetAll(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "successfully get all category")
}

func (h *CategoryHandler) Paginate(c *gin.Context) {
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

func (h *CategoryHandler) FindById(c *gin.Context) {
	id, err := helper.ParseUUID(c)
	if err != nil {
		c.Error(err)
		return
	}
	category, err := h.UseCase.FindById(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, category, "success get category")
}

func (h *CategoryHandler) Store(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	created, err := h.UseCase.Store(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, created, "success create category")
}

func (h *CategoryHandler) Update(c *gin.Context) {
	var req dto.UpdateCategoryRequest
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

	response.OK(c, res, "success update category")
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := helper.ParseUUID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.UseCase.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	response.OK(c, nil, "success delete category")
}

func (h *CategoryHandler) Restore(c *gin.Context) {
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

	response.OK(c, res, "success restore category")
}
