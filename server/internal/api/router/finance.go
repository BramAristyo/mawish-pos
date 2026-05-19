package router

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func COARoutes(r *gin.RouterGroup, h *handler.COAHandler) {
	r.GET("get-all", h.GetAll)
	r.GET("get-all/operational", h.GetAllOperational)
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
}

func CashTransactionRoutes(r *gin.RouterGroup, h *handler.CashTransactionHandler) {
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

func LedgerRoutes(r *gin.RouterGroup, h *handler.LedgerHandler) {
	r.GET("/transactions", h.TransactionList)
	r.GET("/cash-flow", h.CashFlowStatement)
}
