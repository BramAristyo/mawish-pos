package usecase

import (
	"context"
	"fmt"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	Repo       *repository.UserRepository
	LogUseCase *AuditLogUseCase
}

func NewUserUseCase(r *repository.UserRepository, log *AuditLogUseCase) *UserUseCase {
	return &UserUseCase{
		Repo:       r,
		LogUseCase: log,
	}
}

func (u *UserUseCase) GetAll(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := u.Repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		responses = append(responses, dto.ToUserResponse(&users[i]))
	}

	return responses, nil
}

func (u *UserUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	user, err := u.Repo.FindById(id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.ToUserResponse(&user), nil
}

func (u *UserUseCase) Store(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return dto.UserResponse{}, err
	}

	req.Password = string(hashed)

	user := dto.ToCreateUserModel(&req)

	stored, err := u.Repo.Store(&user)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.UserResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Email",
					Tag:      "unique",
					Value:    req.Email,
					Message:  "This email already exists.",
				},
			}
		}
		return dto.UserResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityUser,
		EntityID:    &stored.ID,
		Description: "User created a new user account: " + stored.Email,
	})

	return dto.ToUserResponse(&stored), nil
}

func (u *UserUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)

	user := dto.ToUpdateUserModel(&req)

	updatedUser, err := u.Repo.Update(id, &user)
	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.UserResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "Email",
					Tag:      "unique",
					Value:    req.Email,
					Message:  "This email already exists.",
				},
			}
		}
		return dto.UserResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityUser,
		EntityID:    &updatedUser.ID,
		Description: "User updated user account: " + updatedUser.Email,
	})

	return dto.ToUserResponse(&updatedUser), nil
}

func (u *UserUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userId, _ := helper.ExtractUserID(ctx)

	user, err := u.Repo.FindById(id)
	if err != nil {
		return err
	}

	err = u.Repo.Delete(id)
	if err != nil {
		return err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionDelete,
		Entity:      domain.EntityUser,
		EntityID:    &id,
		Description: fmt.Sprintf("User deleted user account: %s", user.Email),
	})

	return nil
}

func (u *UserUseCase) Restore(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)

	if err := u.Repo.Restore(id); err != nil {
		return dto.UserResponse{}, err
	}

	user, err := u.Repo.FindById(id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionRestore,
		Entity:      domain.EntityUser,
		EntityID:    &id,
		Description: fmt.Sprintf("User restored user account: %s", user.Email),
	})

	return dto.ToUserResponse(&user), nil
}
