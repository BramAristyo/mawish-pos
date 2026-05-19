package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Role      domain.Role `json:"role"`
	Email     string      `json:"email"`
	DeletedAt *string     `json:"deletedAt,omitempty"`
	CreatedAt string      `json:"createdAt"`
}

func ToUserResponse(u *domain.User) UserResponse {
	resp := UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Role:      u.Role,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if u.DeletedAt.Valid {
		at := u.DeletedAt.Time.Format("2006-01-02 15:04:05")
		resp.DeletedAt = &at
	}

	return resp
}

func ToCreateUserModel(req *CreateUserRequest) domain.User {
	return domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToUpdateUserModel(req *UpdateUserRequest) domain.User {
	return domain.User{
		Name:  req.Name,
		Email: req.Email,
	}
}
