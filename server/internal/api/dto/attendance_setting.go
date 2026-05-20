package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
)

type AttendanceSettingResponse struct {
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Radius    float64 `json:"radius"`
	UpdatedAt string  `json:"updatedAt"`
}

type AttendanceSettingRequest struct {
	Lat    float64 `json:"lat" binding:"required"`
	Lng    float64 `json:"lng" binding:"required"`
	Radius float64 `json:"radius" binding:"required"`
}

func ToAttendanceSettingResponse(s domain.AttendanceSetting) AttendanceSettingResponse {
	return AttendanceSettingResponse{
		Lat:       s.Lat,
		Lng:       s.Lng,
		Radius:    s.Radius,
		UpdatedAt: s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToAttendanceSettingDomain(req AttendanceSettingRequest) domain.AttendanceSetting {
	return domain.AttendanceSetting{
		Lat:    req.Lat,
		Lng:    req.Lng,
		Radius: req.Radius,
	}
}
