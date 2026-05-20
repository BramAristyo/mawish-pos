package router

import (
	"fmt"
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/api/middleware"
	"github.com/BramAristyo/mawish-pos/server/internal/dependency"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *dependency.Handlers, cfg *config.Config) {

	r.GET("/slow", func(c *gin.Context) {
		for i := range 5 {
			fmt.Println(i)
			time.Sleep(1 * time.Second)
		}
		c.JSON(200, gin.H{"message": "done"})
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		users := v1.Group("/users", middleware.Authentication(cfg))
		categories := v1.Group("/categories", middleware.Authentication(cfg))
		products := v1.Group("/products", middleware.Authentication(cfg))
		modifierGroups := v1.Group("/modifier-groups", middleware.Authentication(cfg))
		bundling := v1.Group("/bundling", middleware.Authentication(cfg))
		taxes := v1.Group("/taxes", middleware.Authentication(cfg))
		discounts := v1.Group("/discounts", middleware.Authentication(cfg))
		shifts := v1.Group("/shifts", middleware.Authentication(cfg))
		salesTypes := v1.Group("/sales-types", middleware.Authentication(cfg))
		orders := v1.Group("/orders", middleware.Authentication(cfg))
		reports := v1.Group("/reports", middleware.Authentication(cfg))
		dashboard := v1.Group("/dashboard", middleware.Authentication(cfg))
		coa := v1.Group("/coa", middleware.Authentication(cfg))
		employees := v1.Group("/employees", middleware.Authentication(cfg))
		attendances := v1.Group("/attendances", middleware.Authentication(cfg))
		attendanceSettings := v1.Group("/attendance-settings", middleware.Authentication(cfg))
		payrolls := v1.Group("/payrolls", middleware.Authentication(cfg))
		shiftSchedules := v1.Group("/shift-schedules", middleware.Authentication(cfg))
		cashTransactions := v1.Group("/cash-transactions", middleware.Authentication(cfg))
		ledger := v1.Group("/ledger", middleware.Authentication(cfg))

		v1.POST("/", h.Auth.Login)
		v1.GET("/me", middleware.Authentication(cfg), h.Auth.Me)

		UserRoutes(users, h.User)
		CategoryRoutes(categories, h.Category)
		ProductRoutes(products, h.Product)
		ModifierRoutes(modifierGroups, h.ModifierGroup)
		BundlingRoutes(bundling, h.Bundling)
		TaxRoutes(taxes, h.Tax)
		DiscountRoutes(discounts, h.Discount)
		ShiftRoutes(shifts, h.Shift)
		SalesTypeRoutes(salesTypes, h.SalesType)
		OrderRoutes(orders, h.Order)
		ReportRoutes(reports, h.Report)
		DashboardRoutes(dashboard, h.Dashboard)
		COARoutes(coa, h.COA)
		EmployeeRoutes(employees, h.Employee)
		AttendanceRoutes(attendances, h.Attendance)
		AttendanceSettingRoutes(attendanceSettings, h.AttendanceSetting)
		PayrollRoutes(payrolls, h.Payroll)
		ShiftScheduleRoutes(shiftSchedules, h.ShiftSchedule)
		CashTransactionRoutes(cashTransactions, h.CashTransaction)
		LedgerRoutes(ledger, h.Ledger)
	}
}
