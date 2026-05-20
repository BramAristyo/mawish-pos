package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type AttendanceSettingHandler struct {
	usecase *usecase.AttendanceSettingUseCase
}

func NewAttendanceSettingHandler(u *usecase.AttendanceSettingUseCase) *AttendanceSettingHandler {
	return &AttendanceSettingHandler{
		usecase: u,
	}
}

func (h *AttendanceSettingHandler) Get(c *gin.Context) {
	res, err := h.usecase.Get(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success get attendance settings")
}

func (h *AttendanceSettingHandler) Update(c *gin.Context) {
	var req dto.AttendanceSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.usecase.Update(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "success update attendance settings")
}
