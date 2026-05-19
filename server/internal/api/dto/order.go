package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderItemResponse struct {
	ID           uuid.UUID       `json:"id"`
	OrderID      uuid.UUID       `json:"orderId"`
	ProductID    *uuid.UUID      `json:"productId"`
	BundlingID   *uuid.UUID      `json:"bundlingId,omitempty"`
	DiscountID   *uuid.UUID      `json:"discountId,omitempty"`
	ProductName  string          `json:"productName"`
	ProductPrice decimal.Decimal `json:"productPrice"`
	ProductCogs  decimal.Decimal `json:"productCogs"`
	Quantity     int             `json:"quantity"`

	DiscountAmount decimal.Decimal `json:"discountAmount"`
	Subtotal       decimal.Decimal `json:"subtotal"`

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	Product   *ProductResponse            `json:"product,omitempty"`
	Bundling  *BundlingPackageResponse    `json:"bundling,omitempty"`
	Discount  *DiscountResponse           `json:"discount,omitempty"`
	Modifiers []OrderItemModifierResponse `json:"modifiers,omitempty"`
}

type OrderItemModifierResponse struct {
	ID               uuid.UUID       `json:"id"`
	ModifierOptionID *uuid.UUID      `json:"modifierOptionId"`
	ModifierName     string          `json:"modifierName"`
	PriceAdjustment  decimal.Decimal `json:"priceAdjustment"`
	CogsAdjustment   decimal.Decimal `json:"cogsAdjustment"`
}

type PaymentMethodResponse struct {
	ID        uuid.UUID            `json:"id"`
	OrderID   uuid.UUID            `json:"orderId"`
	Method    domain.PaymentMethod `json:"paymentMethod"`
	Amount    decimal.Decimal      `json:"amount"`
	CreatedAt string               `json:"createdAt"`
}

type OrderResponse struct {
	ID          uuid.UUID  `json:"id"`
	ShiftID     uuid.UUID  `json:"shiftId"`
	CashierID   uuid.UUID  `json:"cashierId"`
	SalesTypeID uuid.UUID  `json:"salesTypeId"`
	TaxID       *uuid.UUID `json:"taxId,omitempty"`
	DiscountID  *uuid.UUID `json:"discountId,omitempty"`
	OrderNumber string     `json:"orderNumber"`

	Subtotal       decimal.Decimal `json:"subtotal"`
	DiscountAmount decimal.Decimal `json:"discountAmount"`
	TaxAmount      decimal.Decimal `json:"taxAmount"`
	ChargeAmount   decimal.Decimal `json:"chargeAmount"`
	Total          decimal.Decimal `json:"total"`

	Status     domain.OrderStatus `json:"status"`
	VoidReason *string            `json:"voidReason,omitempty"`
	VoidedBy   *uuid.UUID         `json:"voidedBy,omitempty"`
	VoidedAt   *string            `json:"voidedAt,omitempty"`
	CreatedAt  string             `json:"createdAt"`
	UpdatedAt  string             `json:"updatedAt"`

	Shift        *ShiftResponse    `json:"shift,omitempty"`
	Cashier      UserResponse      `json:"cashier"`
	SalesType    SalesTypeResponse `json:"salesType"`
	Tax          *TaxResponse      `json:"tax,omitempty"`
	Discount     *DiscountResponse `json:"discount,omitempty"`
	VoidedByUser *UserResponse     `json:"voidedByUser,omitempty"`

	Items    []OrderItemResponse     `json:"items,omitempty"`
	Payments []PaymentMethodResponse `json:"payments,omitempty"`
}

type OrderResponsePagination struct {
	Data []OrderResponse `json:"data"`
	Meta filter.Meta     `json:"meta"`
}

type CreateOrderItemRequest struct {
	ProductID         *string  `json:"productId" binding:"required_without=BundlingID,omitempty,uuid"`
	BundlingID        *string  `json:"bundlingId" binding:"required_without=ProductID,omitempty,uuid"`
	DiscountID        *string  `json:"discountId" binding:"omitempty,uuid"`
	ModifierOptionIDs []string `json:"modifierOptionIds" binding:"omitempty,dive,uuid"`
	Quantity          int      `json:"quantity" binding:"required,min=1"`
}

type CreateOrderPaymentRequest struct {
	Method domain.PaymentMethod `json:"method" binding:"required,oneof=cash qris card"`
	Amount decimal.Decimal      `json:"amount" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	SalesTypeID string  `json:"salesTypeId" binding:"required,uuid"`
	TaxID       *string `json:"taxId" binding:"omitempty,uuid"`
	DiscountID  *string `json:"discountId" binding:"omitempty,uuid"`

	Items    []CreateOrderItemRequest    `json:"items" binding:"required,min=1,dive"`
	Payments []CreateOrderPaymentRequest `json:"payments" binding:"required,min=1,dive"`
}

type VoidOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

func ToOrderItemResponse(item domain.OrderItem) OrderItemResponse {
	res := OrderItemResponse{
		ID:             item.ID,
		OrderID:        item.OrderID,
		ProductID:      item.ProductID,
		BundlingID:     item.BundlingID,
		DiscountID:     item.DiscountID,
		ProductName:    item.ProductName,
		ProductPrice:   item.ProductPrice,
		ProductCogs:    item.ProductCogs,
		Quantity:       item.Quantity,
		DiscountAmount: item.DiscountAmount,
		Subtotal:       item.Subtotal,
		CreatedAt:      item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if item.Product != nil {
		p := ToProductResponse(item.Product)
		res.Product = &p
	}
	if item.Bundling != nil {
		b := ToBundlingPackageResponse(item.Bundling)
		res.Bundling = &b
	}
	if item.Discount != nil {
		d := ToDiscountResponse(item.Discount)
		res.Discount = &d
	}

	if len(item.Modifiers) > 0 {
		mods := make([]OrderItemModifierResponse, 0, len(item.Modifiers))
		for _, m := range item.Modifiers {
			mods = append(mods, OrderItemModifierResponse{
				ID:               m.ID,
				ModifierOptionID: m.ModifierOptionID,
				ModifierName:     m.ModifierName,
				PriceAdjustment:  m.PriceAdjustment,
				CogsAdjustment:   m.CogsAdjustment,
			})
		}
		res.Modifiers = mods
	}

	return res
}

func ToPaymentMethodResponse(p domain.Payment) PaymentMethodResponse {
	return PaymentMethodResponse{
		ID:        p.ID,
		OrderID:   p.OrderID,
		Method:    p.Method,
		Amount:    p.Amount,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToOrderResponse(order domain.Order) OrderResponse {
	var voidedAt *string
	if order.VoidedAt != nil {
		vStr := order.VoidedAt.Format("2006-01-02 15:04:05")
		voidedAt = &vStr
	}

	cashierRes := ToUserResponse(&order.Cashier)
	salesTypeRes := ToSalesTypeResponse(&order.SalesType)

	res := OrderResponse{
		ID:             order.ID,
		ShiftID:        order.ShiftID,
		CashierID:      order.CashierID,
		SalesTypeID:    order.SalesTypeID,
		TaxID:          order.TaxID,
		DiscountID:     order.DiscountID,
		OrderNumber:    order.OrderNumber,
		Subtotal:       order.Subtotal,
		DiscountAmount: order.DiscountAmount,
		TaxAmount:      order.TaxAmount,
		ChargeAmount:   order.ChargeAmount,
		Total:          order.Total,
		Status:         order.Status,
		VoidReason:     order.VoidReason,
		VoidedBy:       order.VoidedBy,
		VoidedAt:       voidedAt,
		CreatedAt:      order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      order.UpdatedAt.Format("2006-01-02 15:04:05"),
		Cashier:        cashierRes,
		SalesType:      salesTypeRes,
	}

	// Only include Shift when it's present (avoid always setting pointer which defeats omitempty)
	if order.Shift.ID != uuid.Nil {
		shiftRes := ToShiftResponse(&order.Shift)
		res.Shift = &shiftRes
	}

	if order.Tax != nil {
		t := ToTaxResponse(order.Tax)
		res.Tax = &t
	}
	if order.Discount != nil {
		d := ToDiscountResponse(order.Discount)
		res.Discount = &d
	}
	if order.VoidedByUser != nil {
		v := ToUserResponse(order.VoidedByUser)
		res.VoidedByUser = &v
	}

	items := make([]OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, ToOrderItemResponse(item))
	}
	res.Items = items

	payments := make([]PaymentMethodResponse, 0, len(order.Payments))
	for _, payment := range order.Payments {
		payments = append(payments, ToPaymentMethodResponse(payment))
	}
	res.Payments = payments

	return res
}

func ToOrderResponsePagination(orders []domain.Order, p filter.PaginationWithInputFilter, totalRows int64) OrderResponsePagination {
	data := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		data = append(data, ToOrderResponse(order))
	}

	return OrderResponsePagination{
		Data: data,
		Meta: p.ToMeta(totalRows),
	}
}
