package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmployeeHandler struct {
	usecase *usecase.EmployeeUseCase
}

func NewEmployeeHandler(u *usecase.EmployeeUseCase) *EmployeeHandler {
	return &EmployeeHandler{
		usecase: u,
	}
}

func (h *EmployeeHandler) GetAll(c *gin.Context) {
	res, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "Success fetch all employees")
}

func (h *EmployeeHandler) Paginate(c *gin.Context) {
	var req filter.PaginationWithInputFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.Paginate(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OKPaginate(c, res.Data, res.Meta)
}

func (h *EmployeeHandler) FindById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.FindById(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "Success fetch employee by id")
}

func (h *EmployeeHandler) Store(c *gin.Context) {
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.Store(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, res, "Success create employee")
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.Update(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "Success update employee")
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}

	response.OK(c, nil, "Success delete employee")
}

func (h *EmployeeHandler) Restore(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.Restore(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "Success restore employee")
}
