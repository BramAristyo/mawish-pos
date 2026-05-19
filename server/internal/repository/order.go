package repository

import (
	"context"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/persistence/database"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Paginate(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.Order, error) {
	var totalRows int64

	allowedFields := map[string]string{
		"order_number": "order_number",
		"status":       "status",
		"total":        "total",
		"created_at":   "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.Order{}), req.DynamicFilter, []string{"order_number"}, allowedFields)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	if totalRows == 0 {
		return 0, nil, nil
	}

	orders := make([]domain.Order, 0, req.PaginationInput.PageSize)
	if err := q.
		Preload("Cashier").
		Preload("SalesType").
		Offset(req.Offset()).
		Limit(req.PaginationInput.PageSize).
		Find(&orders).
		Error; err != nil {
		return 0, nil, err
	}

	return totalRows, orders, nil
}

func (r *OrderRepository) FindById(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	var existing domain.Order
	if err := r.DB.WithContext(ctx).
		Where("id = ?", id).
		Preload("Shift").
		Preload("Cashier").
		Preload("SalesType").
		Preload("Tax").
		Preload("Discount").
		Preload("VoidedByUser").
		Preload("Items.Product").
		Preload("Items.Bundling").
		Preload("Items.Discount").
		Preload("Items.Modifiers").
		Preload("Payments").
		First(&existing).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Order{}, usecase_errors.NotFound
		}
		return domain.Order{}, err
	}

	return existing, nil
}

func (r *OrderRepository) Store(ctx context.Context, order *domain.Order) (domain.Order, error) {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return domain.Order{}, err
	}

	return r.FindById(ctx, order.ID)
}

func (r *OrderRepository) Void(ctx context.Context, order *domain.Order) (domain.Order, error) {
	result := r.DB.WithContext(ctx).Model(order).Updates(map[string]any{
		"status":      domain.OrderVoided,
		"void_reason": order.VoidReason,
		"voided_by":   order.VoidedBy,
		"voided_at":   order.VoidedAt,
	})

	if result.Error != nil {
		return domain.Order{}, result.Error
	}

	if result.RowsAffected == 0 {
		return domain.Order{}, usecase_errors.NotFound
	}

	return r.FindById(ctx, order.ID)
}

func (r *OrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (domain.Order, error) {
	var existing domain.Order
	if err := r.DB.WithContext(ctx).
		Where("order_number = ?", orderNumber).
		Preload("Shift").
		Preload("Cashier").
		Preload("SalesType").
		Preload("Tax").
		Preload("Discount").
		Preload("VoidedByUser").
		Preload("Items.Product").
		Preload("Items.Bundling").
		Preload("Items.Discount").
		Preload("Items.Modifiers").
		Preload("Payments").
		First(&existing).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Order{}, usecase_errors.NotFound
		}
		return domain.Order{}, err
	}

	return existing, nil
}

func (r *OrderRepository) GetLatestOrder(ctx context.Context) (domain.Order, error) {
	var latest domain.Order
	if err := r.DB.WithContext(ctx).
		Order("created_at desc").
		First(&latest).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Order{}, usecase_errors.NotFound
		}
		return domain.Order{}, err
	}

	return latest, nil
}

func (r *OrderRepository) SalesSummary(ctx context.Context, req filter.DynamicFilter) (domain.SalesSummary, error) {
	var summary domain.SalesSummary

	allowedFields := map[string]string{
		"created_at": "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.Order{}), req, nil, allowedFields)

	err := q.
		Select(`
			COAlESCE(SUM(subtotal), 0) AS gross_sales,
			COALESCE(SUM(discount_amount), 0) AS discounts,
			COALESCE(SUM(subtotal - discount_amount), 0) AS net_sales,
			COALESCE(SUM(charge_amount), 0) AS gratuity,
			COALESCE(SUM(tax_amount), 0) AS tax,
			COALESCE(SUM(total), 0) AS total
		`).
		Where("status = ?", domain.OrderCompleted).
		Scan(&summary).Error

	if err != nil {
		return domain.SalesSummary{}, err
	}

	return summary, nil
}

func (r *OrderRepository) GrossProfit(ctx context.Context, req filter.DynamicFilter) (domain.GrossProfit, error) {
	var summary domain.GrossProfit

	allowedFields := map[string]string{
		"created_at": "created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.Order{}), req, nil, allowedFields)

	err := q.
		Select(`
			COAlESCE(SUM(subtotal), 0) AS gross_sales,
			COALESCE(SUM(discount_amount), 0) AS discounts,
			COALESCE(SUM(subtotal - discount_amount), 0) AS net_sales,
			COALESCE(SUM(item_agg.total_cogs), 0) AS cogs
		`).
		Joins(`LEFT JOIN(
			SELECT order_id, SUM(product_cogs * quantity) AS total_cogs
			FROM order_items
			GROUP BY order_id
			) item_agg ON item_agg.order_id = orders.id`).
		Where("status = ?", domain.OrderCompleted).
		Scan(&summary).Error

	if err != nil {
		return domain.GrossProfit{}, err
	}

	return summary, nil
}

func (r *OrderRepository) TransactionReport(ctx context.Context, req filter.PaginationWithInputFilter) (int64, []domain.Transaction, error) {
	var totalRows int64

	allowedFields := map[string]string{
		"created_at":   "orders.created_at",
		"order_number": "orders.order_number",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Model(&domain.OrderItem{}), req.DynamicFilter, []string{"orders.order_number"}, allowedFields)

	q = q.
		Select("orders.order_number, TO_CHAR(orders.created_at, 'YYYY-MM-DD HH24:MI:SS') as time, order_items.product_name as product, order_items.subtotal as price").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.status = ?", domain.OrderCompleted)

	if err := q.Count(&totalRows).Error; err != nil {
		return 0, nil, err
	}

	if totalRows == 0 {
		return 0, nil, nil
	}

	transactions := make([]domain.Transaction, 0, req.PaginationInput.PageSize)
	err := q.Offset(req.Offset()).
		Limit(req.PaginationInput.PageSize).
		Scan(&transactions).Error

	if err != nil {
		return 0, nil, err
	}

	return totalRows, transactions, nil
}

func (r *OrderRepository) SalesSummaryDashboard(ctx context.Context, req filter.DynamicFilter) (domain.SalesSummaryDashboard, error) {
	var summary struct {
		GrossSales       decimal.Decimal
		NetSales         decimal.Decimal
		TransactionCount int64
		Cogs             decimal.Decimal
	}

	allowedFields := map[string]string{
		"created_at": "orders.created_at",
	}

	q := database.BuildQuery(r.DB.WithContext(ctx).Table("orders"), req, nil, allowedFields)

	err := q.
		Select(`
			COALESCE(SUM(orders.subtotal), 0) AS gross_sales,
			COALESCE(SUM(orders.subtotal - orders.discount_amount), 0) AS net_sales,
			COUNT(orders.id) AS transaction_count,
			COALESCE(SUM(item_agg.total_cogs), 0) AS cogs
		`).
		Joins(`LEFT JOIN (
			SELECT order_id, SUM(product_cogs * quantity) AS total_cogs
			FROM order_items
			GROUP BY order_id
		) item_agg ON item_agg.order_id = orders.id`).
		Where("status = ?", domain.OrderCompleted).
		Scan(&summary).Error

	if err != nil {
		return domain.SalesSummaryDashboard{}, err
	}

	result := domain.SalesSummaryDashboard{
		GrossSales:       summary.GrossSales,
		NetSales:         summary.NetSales,
		GrossProfit:      summary.NetSales.Sub(summary.Cogs),
		TransactionCount: summary.TransactionCount,
	}
	result.Calculate()

	return result, nil
}
