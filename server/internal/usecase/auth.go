package usecase

import (
	"context"
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/constant"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/config"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	Repo       *repository.UserRepository
	Cfg        *config.Config
	LogUseCase *AuditLogUseCase
}

func NewAuthUseCase(repo *repository.UserRepository, cfg *config.Config, log *AuditLogUseCase) *AuthUseCase {
	return &AuthUseCase{
		Repo:       repo,
		Cfg:        cfg,
		LogUseCase: log,
	}
}

func (u *AuthUseCase) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := u.Repo.FindByEmail(req.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	if user.DeletedAt.Valid {
		return dto.LoginResponse{}, usecase_errors.UserNotActive
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.LoginResponse{}, usecase_errors.InvalidPassword
	}

	atc := jwt.MapClaims{}

	atc["userID"] = user.ID
	atc["email"] = user.Email
	atc["role"] = user.Role
	atc["exp"] = time.Now().Add(u.Cfg.JWT.AccessTokenExpireDuration).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)
	token, err := at.SignedString([]byte(u.Cfg.JWT.Secret))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      user.ID,
		Action:      domain.ActionLogin,
		Entity:      domain.EntityUser,
		EntityID:    &user.ID,
		Description: "User logged in successfully",
	})

	return dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(&user),
	}, nil
}

func (u *AuthUseCase) Me(ctx context.Context) (dto.UserResponse, error) {
	userIDStr, ok := ctx.Value(constant.CtxUserID).(string)
	if !ok {
		return dto.UserResponse{}, usecase_errors.TokenRequired
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return dto.UserResponse{}, usecase_errors.TokenInvalid
	}

	user, err := u.Repo.FindById(userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.ToUserResponse(&user), nil
}
