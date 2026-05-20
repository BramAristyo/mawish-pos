package domain

import "time"

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
