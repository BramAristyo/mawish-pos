package router

import (
	"github.com/BramAristyo/mawish-pos/server/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func EmployeeRoutes(r *gin.RouterGroup, h *handler.EmployeeHandler) {
	r.GET("", h.Paginate)
	r.GET("/all", h.GetAll)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
}

func AttendanceRoutes(r *gin.RouterGroup, h *handler.AttendanceHandler) {
	r.GET("", h.Paginate)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.PUT("/:id/photo", h.ConfirmImage)
}

func AttendanceSettingRoutes(r *gin.RouterGroup, h *handler.AttendanceSettingHandler) {
	r.GET("", h.Get)
	r.POST("", h.Update)
}

func ShiftScheduleRoutes(r *gin.RouterGroup, h *handler.ShiftScheduleHandler) {
	r.GET("", h.Paginate)
	r.GET("/all", h.GetAll)
	r.GET("/:id", h.FindById)
	r.POST("", h.Store)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.PATCH("/:id/restore", h.Restore)
}

func PayrollRoutes(r *gin.RouterGroup, h *handler.PayrollHandler) {
	r.GET("", h.Paginate)
	r.POST("", h.Store)
}
