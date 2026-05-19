package seeder

import (
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	// Refactored with hardcoded UUIDs and deterministic logic
	SeedUserData(db)
	SeedCOAData(db)

	SeedCategoryData(db)
	SeedModifierGroupData(db)
	SeedTaxData(db)
	SeedDiscountData(db)
	SeedSalesTypeData(db)

	SeedProductData(db)
	SeedProductModifierData(db)
	SeedBundlingData(db)

	SeedEmployeeData(db)
	SeedShiftScheduleData(db)
	SeedAttendanceData(db)
	SeedPayrollData(db)

	SeedShiftData(db)
	SeedOrderData(db)

	SeedCashTransactionData(db)
	SeedLedgerData(db)
}

