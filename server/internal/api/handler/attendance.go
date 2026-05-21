package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttendanceHandler struct {
	usecase *usecase.AttendanceUseCase
}

func NewAttendanceHandler(u *usecase.AttendanceUseCase) *AttendanceHandler {
	return &AttendanceHandler{
		usecase: u,
	}
}

func (h *AttendanceHandler) Paginate(c *gin.Context) {
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

func (h *AttendanceHandler) Store(c *gin.Context) {
	var req dto.CreateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.Store(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, res, "success create attendance")
}

func (h *AttendanceHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.Update(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success update attendance")
}

func (h *AttendanceHandler) FindById(c *gin.Context) {
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

	response.OK(c, res, "success get attendance")
}

func (h *AttendanceHandler) ConfirmImage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	var req dto.ConfirmImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.ConfirmAttendanceImage(c.Request.Context(), id, req.Key)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success confirm attendance image")
}
