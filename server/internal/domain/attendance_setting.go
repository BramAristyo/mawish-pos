package domain

import (
	"math"
	"time"
)

type AttendanceSetting struct {
	ID        uint      `gorm:"primaryKey"`
	Lat       float64   `gorm:"not null"`
	Lng       float64   `gorm:"not null"`
	Radius    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (AttendanceSetting) TableName() string {
	return "attendance_settings"
}

func (as *AttendanceSetting) IsValidDistance(lat float64, lng float64) bool {
	const R = 6371000

	dLat := (lat - as.Lat) * (math.Pi / 180.0)
	dLng := (lng - as.Lng) * (math.Pi / 180.0)

	lat1 := as.Lat * (math.Pi / 180.0)
	lat2 := lat * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLng/2)*math.Sin(dLng/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	dist := R * c

	return dist <= as.Radius
}
