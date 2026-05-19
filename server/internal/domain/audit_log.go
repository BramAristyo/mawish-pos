package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditAction string
type AuditEntity string

const (
	ActionCreate  AuditAction = "CREATE"
	ActionUpdate  AuditAction = "UPDATE"
	ActionDelete  AuditAction = "DELETE"
	ActionRestore AuditAction = "RESTORE"
	ActionLogin   AuditAction = "LOGIN"
	ActionVoid    AuditAction = "VOID"
)

const (
	EntityProduct        AuditEntity = "products"
	EntityCategory       AuditEntity = "categories"
	EntityUser           AuditEntity = "users"
	EntityModifierGroup  AuditEntity = "modifier_groups"
	EntityModifierOption AuditEntity = "modifier_options"
	EntityShift          AuditEntity = "shifts"
	EntityOrder          AuditEntity = "orders"
	EntityTax            AuditEntity = "taxes"
	EntityBundling       AuditEntity = "bundling"
	EntityDiscount       AuditEntity = "discounts"
	EntitySalesType      AuditEntity = "sales_types"
	EntityLedger         AuditEntity = "ledgers"
	EntityEmployee       AuditEntity = "employees"
)

type AuditLog struct {
	ID          uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID
	Action      AuditAction
	Entity      AuditEntity
	EntityID    *uuid.UUID
	Description string

	CreatedAt time.Time `gorm:"autoCreateTime"`

	User *User `gorm:"foreignKey:UserID"`
}
