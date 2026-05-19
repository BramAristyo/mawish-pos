package usecase

import (
	"context"
	"fmt"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/google/uuid"
)

type COAUseCase struct {
	Repo       *repository.COARepository
	LogUseCase *AuditLogUseCase
}

func NewCOAUseCase(repo *repository.COARepository, log *AuditLogUseCase) *COAUseCase {
	return &COAUseCase{
		Repo:       repo,
		LogUseCase: log,
	}
}

func (u *COAUseCase) GetAll(ctx context.Context) ([]dto.ChartOfAccountResponse, error) {
	coas, err := u.Repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return dto.ToCOAResponses(coas), nil
}

func (u *COAUseCase) GetAllOperational(ctx context.Context) ([]dto.ChartOfAccountResponse, error) {
	coas, err := u.Repo.GetlAllOperational(ctx)
	if err != nil {
		return nil, err
	}
	return dto.ToCOAResponses(coas), nil
}

func (u *COAUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.ChartOfAccountResponsePagination, error) {
	totalRows, coas, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.ChartOfAccountResponsePagination{}, err
	}

	res := dto.ToCOAResponses(coas)
	return dto.ToCOAResponsePagination(res, req, totalRows), nil
}

func (u *COAUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.ChartOfAccountResponse, error) {
	coa, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.ChartOfAccountResponse{}, err
	}
	return dto.ToCOAResponse(&coa), nil
}

func (u *COAUseCase) Store(ctx context.Context, req dto.CreateCOARequest) (dto.ChartOfAccountResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)
	coa := dto.ToCOAModel(&req)

	stored, err := u.Repo.Store(ctx, &coa)
	if err != nil {
		return dto.ChartOfAccountResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityLedger, // Using EntityLedger as placeholder or add EntityCOA
		EntityID:    &stored.ID,
		Description: fmt.Sprintf("User created a new COA: %s (%s)", stored.Name, stored.Type),
	})

	return dto.ToCOAResponse(&stored), nil
}

func (u *COAUseCase) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCOARequest) (dto.ChartOfAccountResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)
	coa := dto.ToUpdateCOAModel(&req)

	updated, err := u.Repo.Update(ctx, id, &coa)
	if err != nil {
		return dto.ChartOfAccountResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityLedger,
		EntityID:    &id,
		Description: fmt.Sprintf("User updated COA: %s (%s)", updated.Name, updated.Type),
	})

	return dto.ToCOAResponse(&updated), nil
}

func (u *COAUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userId, _ := helper.ExtractUserID(ctx)

	coa, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return err
	}

	if err := u.Repo.Delete(ctx, id); err != nil {
		return err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionDelete,
		Entity:      domain.EntityLedger,
		EntityID:    &id,
		Description: fmt.Sprintf("User deleted COA: %s", coa.Name),
	})

	return nil
}

func (u *COAUseCase) Restore(ctx context.Context, id uuid.UUID) (dto.ChartOfAccountResponse, error) {
	userId, _ := helper.ExtractUserID(ctx)

	if err := u.Repo.Restore(ctx, id); err != nil {
		return dto.ChartOfAccountResponse{}, err
	}

	coa, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.ChartOfAccountResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionRestore,
		Entity:      domain.EntityLedger,
		EntityID:    &id,
		Description: fmt.Sprintf("User restored COA: %s", coa.Name),
	})

	return dto.ToCOAResponse(&coa), nil
}
