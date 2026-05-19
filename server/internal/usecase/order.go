package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/api/dto"
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/repository"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/helper"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderUseCase struct {
	Repo               *repository.OrderRepository
	ShiftRepo          *repository.ShiftRepository
	SalesType          *repository.SalesTypeRepository
	TaxRepo            *repository.TaxRepository
	DiscountRepo       *repository.DiscountRepository
	ProductRepo        *repository.ProductRepository
	BundlingRepo       *repository.BundlingRepository
	ModifierOptionRepo *repository.ModifierOptionRepository
	LogUseCase         *AuditLogUseCase
}

func NewOrderUseCase(
	r *repository.OrderRepository,
	shift *repository.ShiftRepository,
	salesType *repository.SalesTypeRepository,
	tax *repository.TaxRepository,
	discount *repository.DiscountRepository,
	product *repository.ProductRepository,
	bundling *repository.BundlingRepository,
	modifierOption *repository.ModifierOptionRepository,
	log *AuditLogUseCase,
) *OrderUseCase {
	return &OrderUseCase{
		Repo:               r,
		ShiftRepo:          shift,
		SalesType:          salesType,
		TaxRepo:            tax,
		DiscountRepo:       discount,
		ProductRepo:        product,
		BundlingRepo:       bundling,
		ModifierOptionRepo: modifierOption,
		LogUseCase:         log,
	}
}

func (u *OrderUseCase) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (dto.OrderResponsePagination, error) {
	totalRows, orders, err := u.Repo.Paginate(ctx, req)
	if err != nil {
		return dto.OrderResponsePagination{}, err
	}

	return dto.ToOrderResponsePagination(orders, req, totalRows), nil
}

func (u *OrderUseCase) FindById(ctx context.Context, id uuid.UUID) (dto.OrderResponse, error) {
	order, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	return dto.ToOrderResponse(order), nil
}

// CalculateAll builds the order the same way as Store (including calculations) but does NOT persist it to the DB.
// It returns the calculated order response for preview purposes.
func (u *OrderUseCase) CalculateAll(ctx context.Context, req dto.CreateOrderRequest) (dto.OrderResponse, error) {
	order, err := u.buildOrder(ctx, req)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	// Do not include populated Shift in the preview response.
	// Zero out the Shift so dto.ToOrderResponse omits it (omitempty).
	order.Shift = domain.Shift{}

	return dto.ToOrderResponse(*order), nil
}

// Store builds the order and persists it. The order-building logic is shared with CalculateAll via buildOrder.
func (u *OrderUseCase) Store(ctx context.Context, req dto.CreateOrderRequest) (dto.OrderResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	order, err := u.buildOrder(ctx, req)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	var created domain.Order

	for range 3 {
		created, err = u.Repo.Store(ctx, order)
		if err != nil {
			if usecase_errors.IsUniqueViolation(err) {
				newNum, genErr := u.generateOrderNumber(ctx)
				if genErr != nil {
					return dto.OrderResponse{}, genErr
				}
				order.OrderNumber = newNum
				continue
			}

			return dto.OrderResponse{}, err
		}

		break
	}

	if err != nil {
		if usecase_errors.IsUniqueViolation(err) {
			return dto.OrderResponse{}, &usecase_errors.CustomFieldErrors{
				{
					Property: "OrderNumber",
					Tag:      "unique",
					Value:    order.OrderNumber,
					Message:  "This order number already exists.",
				},
			}
		}
		return dto.OrderResponse{}, err
	}

	u.LogUseCase.Log(ctx, domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionCreate,
		Entity:      domain.EntityOrder,
		EntityID:    &created.ID,
		Description: fmt.Sprintf("Create Order with number: %s", created.OrderNumber),
	})

	return dto.ToOrderResponse(created), nil
}

// buildOrder contains the common logic used by both Store and CalculateAll.
// It builds the domain.Order, fills relations, runs CalculateAll and returns the order pointer.
func (u *OrderUseCase) buildOrder(ctx context.Context, req dto.CreateOrderRequest) (*domain.Order, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return nil, err
	}

	shift, err := u.ShiftRepo.FindOpenShiftByUserId(ctx, userId)
	if err != nil {
		return nil, usecase_errors.NoOpenShift
	}

	salesTypeID, err := uuid.Parse(req.SalesTypeID)
	if err != nil {
		return nil, usecase_errors.InvalidID
	}

	salesType, err := u.SalesType.FindById(ctx, salesTypeID)
	if err != nil {
		return nil, err
	}

	var taxIDPtr *uuid.UUID
	var taxPtr *domain.Tax
	if req.TaxID != nil {
		taxID, err := uuid.Parse(*req.TaxID)
		if err != nil {
			return nil, err
		}
		tax, err := u.TaxRepo.FindById(ctx, taxID)
		if err != nil {
			return nil, err
		}
		taxIDPtr = &taxID
		taxPtr = &tax
	}

	var discountIDPtr *uuid.UUID
	var discountPtr *domain.Discount
	if req.DiscountID != nil {
		discountID, err := uuid.Parse(*req.DiscountID)
		if err != nil {
			return nil, err
		}
		discount, err := u.DiscountRepo.FindById(ctx, discountID)
		if err != nil {
			return nil, err
		}
		discountIDPtr = &discountID
		discountPtr = &discount
	}

	if len(req.Items) == 0 {
		return nil, usecase_errors.EmptyOrderItems
	}

	items := make([]domain.OrderItem, len(req.Items))
	for i, oi := range req.Items {
		var bundlingID uuid.UUID
		var bundling domain.BundlingPackage
		var productID uuid.UUID
		var product domain.Product
		var discountItemID uuid.UUID
		var discountItem domain.Discount

		if oi.BundlingID != nil {
			bundlingID, err = uuid.Parse(*oi.BundlingID)
			if err != nil {
				return nil, err
			}

			bundling, err = u.BundlingRepo.FindById(ctx, bundlingID)
			if err != nil {
				return nil, err
			}

			items[i] = domain.OrderItem{
				BundlingID:   &bundlingID,
				Bundling:     &bundling,
				ProductName:  bundling.Name,
				ProductPrice: bundling.Price,
				ProductCogs:  bundling.Cogs,
				Quantity:     oi.Quantity,
			}
		} else if oi.ProductID != nil {
			productID, err = uuid.Parse(*oi.ProductID)
			if err != nil {
				return nil, err
			}

			product, err = u.ProductRepo.FindById(ctx, productID)
			if err != nil {
				return nil, err
			}

			items[i] = domain.OrderItem{
				ProductID:    &productID,
				Product:      &product,
				ProductName:  product.Name,
				ProductPrice: product.Price,
				ProductCogs:  product.Cogs,
				Quantity:     oi.Quantity,
			}
		} else {
			return nil, usecase_errors.InvalidOrderItem
		}

		if oi.DiscountID != nil {
			discountItemID, err = uuid.Parse(*oi.DiscountID)
			if err != nil {
				return nil, err
			}

			discountItem, err = u.DiscountRepo.FindById(ctx, discountItemID)
			if err != nil {
				return nil, err
			}

			items[i].Discount = &discountItem
		}

		if len(oi.ModifierOptionIDs) > 0 {
			modifiers := make([]domain.OrderItemModifier, 0, len(oi.ModifierOptionIDs))
			for _, mid := range oi.ModifierOptionIDs {
				modUUID, err := uuid.Parse(mid)
				if err != nil {
					return nil, usecase_errors.InvalidID
				}

				option, err := u.ModifierOptionRepo.FindById(ctx, modUUID)
				if err != nil {
					return nil, err
				}

				modifiers = append(modifiers, domain.OrderItemModifier{
					ModifierOptionID: &option.ID,
					ModifierName:     option.Name,
					PriceAdjustment:  option.PriceAdjustment,
					CogsAdjustment:   option.CogsAdjustment,
				})
			}
			items[i].Modifiers = modifiers
		}
	}

	payments := make([]domain.Payment, len(req.Payments))
	for i, p := range req.Payments {
		payments[i] = domain.Payment{
			Method: p.Method,
			Amount: p.Amount,
		}
	}

	orderNumber, err := u.generateOrderNumber(ctx)
	if err != nil {
		return nil, err
	}

	order := &domain.Order{
		ShiftID:     shift.ID,
		CashierID:   userId,
		SalesTypeID: salesTypeID,
		TaxID:       taxIDPtr,
		DiscountID:  discountIDPtr,
		OrderNumber: orderNumber,

		Status:    domain.OrderCompleted,
		Shift:     shift,
		SalesType: salesType,
		Tax:       taxPtr,
		Discount:  discountPtr,

		Items:    items,
		Payments: payments,
	}

	order.CalculateAll()

	return order, nil
}

func (u *OrderUseCase) Void(ctx context.Context, id uuid.UUID, req dto.VoidOrderRequest) (dto.OrderResponse, error) {
	userId, err := helper.ExtractUserID(ctx)
	if err != nil {
		return dto.OrderResponse{}, usecase_errors.TokenInvalid
	}

	order, err := u.Repo.FindById(ctx, id)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	if order.Status == domain.OrderVoided {
		return dto.OrderResponse{}, usecase_errors.OrderAlreadyVoided
	}

	now := time.Now()
	order.Status = domain.OrderVoided
	order.VoidReason = &req.Reason
	order.VoidedBy = &userId
	order.VoidedAt = &now

	voided, err := u.Repo.Void(ctx, &order)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	u.LogUseCase.Log(ctx, domain.AuditLog{
		UserID:      userId,
		Action:      domain.ActionVoid,
		Entity:      domain.EntityOrder,
		EntityID:    &voided.ID,
		Description: fmt.Sprintf("Void Order with number: %s. Reason: %s", voided.OrderNumber, req.Reason),
	})

	return dto.ToOrderResponse(voided), nil
}

func (u *OrderUseCase) generateOrderNumber(ctx context.Context) (string, error) {
	order, err := u.Repo.GetLatestOrder(ctx)
	dateStr := time.Now().Format("20060102")
	prefix := fmt.Sprintf("MW/%s/", dateStr)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return prefix + "00001", nil
		}
		return "", err
	}

	latestNumber := order.OrderNumber
	parts := strings.Split(latestNumber, "/")

	if len(parts) != 3 {
		return prefix + "00001", nil
	}

	lastOrderDate := parts[1]
	if lastOrderDate != dateStr {
		return prefix + "00001", nil
	}

	lastCounter, err := strconv.Atoi(parts[2])
	if err != nil {
		return prefix + "00001", nil
	}

	newCounter := lastCounter + 1

	return fmt.Sprintf("%s/%s/%05d", parts[0], parts[1], newCounter), nil
}
