package usecase

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
)

type CashTransactionUseCase struct {
	Repo       *repository.CashTransactionRepository
	LedgerRepo *repository.LedgerRepository
	LogUseCase *AuditLogUseCase
}

func NewCashTransactionUseCase(
	repo *repository.CashTransactionRepository,
	ledgerRepo *repository.LedgerRepository,
	log *AuditLogUseCase,
) *CashTransactionUseCase {
	return &CashTransactionUseCase{
		Repo:       repo,
		LedgerRepo: ledgerRepo,
		LogUseCase: log,
	}
}

func (u *CashTransactionUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.CashTransactionResponsePagination, error) {
	totalRows, cts, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.CashTransactionResponsePagination{}, err
	}

	return dto.ToCashTransactionPagination(cts, req, totalRows), nil
}

func (u *CashTransactionUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.CashTransactionResponse, error) {
	ct, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.CashTransactionResponse{}, err
	}

	return dto.ToCashTransactionResponse(&ct), nil
}

func (u *CashTransactionUseCase) Store(ctx context.Context, req dto.CashTransactionRequest) (dto.CashTransactionResponse, error) {
	ct, err := dto.ToCashTransactionModel(&req)
	if err != nil {
		return dto.CashTransactionResponse{}, err
	}

	created, err := u.Repo.Store(ctx, &ct)
	if err != nil {
		return dto.CashTransactionResponse{}, err
	}

	_, err = u.LedgerRepo.Store(ctx, created.ToLedgerModel())
	if err != nil {
		return dto.CashTransactionResponse{}, usecase_errors.LedgerRecordFailed
	}

	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.CashTransactionResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityLedger,
		EntityID:    &created.ID,
		Description: "User created a new cash transaction: " + created.COA.Name + " " + created.Amount.String(),
	})

	return dto.ToCashTransactionResponse(&created), nil
}

func (u *CashTransactionUseCase) Update(ctx context.Context, id uuid.UUID, req dto.CashTransactionRequest) (dto.CashTransactionResponse, error) {
	ct, err := dto.ToCashTransactionModel(&req)
	if err != nil {
		return dto.CashTransactionResponse{}, err
	}

	updated, err := u.Repo.Update(ctx, id, &ct)
	if err != nil {
		return dto.CashTransactionResponse{}, err
	}

	err = u.LedgerRepo.Delete(ctx, id, domain.LedgerCashTransaction)
	if err != nil {
		return dto.CashTransactionResponse{}, usecase_errors.LedgerRecordFailed
	}

	_, err = u.LedgerRepo.Store(ctx, updated.ToLedgerModel())
	if err != nil {
		return dto.CashTransactionResponse{}, usecase_errors.LedgerRecordFailed
	}

	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.CashTransactionResponse{}, err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionUpdate,
		Entity:      domain.EntityLedger,
		EntityID:    &updated.ID,
		Description: "User updated a new cash transaction: " + updated.COA.Name + " " + updated.Amount.String(),
	})

	return dto.ToCashTransactionResponse(&updated), nil
}

func (u *CashTransactionUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.Repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = u.LedgerRepo.Delete(ctx, id, domain.LedgerCashTransaction)
	if err != nil {
		return usecase_errors.LedgerRecordFailed
	}

	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return err
	}

	go u.LogUseCase.Log(context.Background(), domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionDelete,
		Entity:      domain.EntityLedger,
		EntityID:    &id,
		Description: "User deleted a cash transaction",
	})

	return nil
}
