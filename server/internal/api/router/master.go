package router

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func CategoryRoutes(r *gin.RouterGroup, h *handler.CategoryHandler) {
	r.GET("get-all", h.GetAll)
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
}

func ProductRoutes(r *gin.RouterGroup, h *handler.ProductHandler) {
	r.GET("get-all", h.GetAll)
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
}

func ModifierRoutes(r *gin.RouterGroup, h *handler.ModifierGroupHandler) {
	r.GET("get-all", h.GetAll)
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
}

func BundlingRoutes(r *gin.RouterGroup, h *handler.BundlingHandler) {
	r.GET("get-all", h.GetAll)
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
}

func TaxRoutes(r *gin.RouterGroup, h *handler.TaxHandler) {
	r.GET("get-all", h.GetAll)
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
	r.PATCH("/:id/activate", h.Activate)
	r.PATCH("/:id/deactivate", h.Deactivate)
}

func DiscountRoutes(r *gin.RouterGroup, h *handler.DiscountHandler) {
	r.GET("get-all", h.GetAll)
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
}

func SalesTypeRoutes(r *gin.RouterGroup, h *handler.SalesTypeHandler) {
	r.GET("get-all", h.GetAll)
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
}
