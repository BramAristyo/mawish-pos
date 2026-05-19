package handler

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/usecase"
	"github.com/BramAristyo/mawish-pos/server/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UseCase *usecase.AuthUseCase
}

func NewAuthHandler(u *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		UseCase: u,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.UseCase.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "login successfully")
}

func (h *AuthHandler) Me(c *gin.Context) {
	res, err := h.UseCase.Me(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, res, "get current user successfully")
}
