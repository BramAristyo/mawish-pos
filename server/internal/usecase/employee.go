package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeUseCase struct {
	Repo       *repository.EmployeeRepository
	LogUseCase *AuditLogUseCase
}

func NewEmployeeUseCase(repo *repository.EmployeeRepository, log *AuditLogUseCase) *EmployeeUseCase {
	return &EmployeeUseCase{
		Repo:       repo,
		LogUseCase: log,
	}
}

func (u *EmployeeUseCase) generateCode(ctx context.Context) string {
	last, err := u.Repo.GetLast(ctx)
	if err != nil {
		return "EMP-001"
	}


	// Format: EMP-001
	parts := strings.Split(last.Code, "-")
	if len(parts) < 2 {
		return "EMP-001"
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "EMP-001"
	}

	return fmt.Sprintf("EMP-%03d", num+1)
}

func (u *EmployeeUseCase) GetAll(ctx context.Context) ([]dto.EmployeeResponse, error) {
	employees, err := u.Repo.GetAll(ctx)
	if err != nil {
		return []dto.EmployeeResponse{}, err
	}

	res := make([]dto.EmployeeResponse, 0, len(employees))
	for _, e := range employees {
		res = append(res, dto.ToEmployeeResponse(&e))
	}

	return res, nil
}

func (u *EmployeeUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.EmployeeResponsePagination, error) {
	totalRows, employees, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.EmployeeResponsePagination{}, err
	}

	res := make([]dto.EmployeeResponse, 0, len(employees))
	for _, e := range employees {
		res = append(res, dto.ToEmployeeResponse(&e))
	}

	return dto.ToEmployeeResponsePagination(res, req, totalRows), nil
}

func (u *EmployeeUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.EmployeeResponse, error) {
	employee, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	return dto.ToEmployeeResponse(&employee), nil
}

func (u *EmployeeUseCase) Store(ctx context.Context, req dto.CreateEmployeeRequest) (dto.EmployeeResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	employee := dto.ToCreateEmployeeModel(&req)
	employee.Code = u.generateCode(ctx)

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Pin), 10)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}
	employee.PinHash = string(hashed)
	employee.HasChangedPIN = false

	stored, err := u.Repo.Store(ctx, &employee)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityEmployee,
		EntityID:    &stored.ID,
		Description: "User created a new employee: " + stored.Name + " (" + stored.Code + ")",
	})

	return dto.ToEmployeeResponse(&stored), nil
}

func (u *EmployeeUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateEmployeeRequest) (dto.EmployeeResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	employee := dto.ToUpdateEmployeeModel(&req)

	if req.Pin != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Pin), 10)
		if err != nil {
			return dto.EmployeeResponse{}, err
		}
		employee.PinHash = string(hashed)
		employee.HasChangedPIN = true
	}

	updated, err := u.Repo.Update(ctx, id, &employee)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityEmployee,
		EntityID:    &updated.ID,
		Description: "User updated employee: " + updated.Name,
	})

	return dto.ToEmployeeResponse(&updated), nil
}

func (u *EmployeeUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return err
	}

	employee, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Repo.Delete(ctx, id); err != nil {
		return err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionDelete,
		Entity:      domain.EntityEmployee,
		EntityID:    &id,
		Description: "User deleted employee: " + employee.Name,
	})

	return nil
}

func (u *EmployeeUseCase) Restore(ctx context.Context, id uuid.UUID) (dto.EmployeeResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	if err := u.Repo.Restore(ctx, id); err != nil {
		return dto.EmployeeResponse{}, err
	}

	employee, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.EmployeeResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionRestore,
		Entity:      domain.EntityEmployee,
		EntityID:    &id,
		Description: "User restored employee: " + employee.Name,
	})

	return dto.ToEmployeeResponse(&employee), nil
}
