package dto

import (
	"fmt"
	"time"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AttendanceResponse struct {
	ID                uuid.UUID       `json:"id"`
	EmployeeCode      string          `json:"employeeCode"`
	EmployeeName      string          `json:"employeeName"`
	Date              string          `json:"date"`
	CheckIn           string          `json:"checkIn"`
	CheckOut          *string         `json:"checkOut"`
	TotalWorkMinutes  *string         `json:"totalWorkMinutes"`
	ShiftScheduleName string          `json:"shiftScheduleName"`
	Notes             *string         `json:"notes"`
	LateMinutes       int             `json:"lateMinutes"`
	DeductionAmount   decimal.Decimal `json:"deductionAmount"`
	DeletedAt         *string         `json:"deletedAt,omitempty"`
	PhotoKey          *string         `json:"photoKey,omitempty"`

	UploadURL string `json:"uploadUrl,omitempty"`
	PhotoURL  string `json:"photoUrl,omitempty"`
}

type AttendanceResponsePagination struct {
	Data []AttendanceResponse `json:"data"`
	Meta filter.Meta          `json:"meta"`
}

type AttendanceRequest struct {
	EmployeeID      string  `json:"employeeId" binding:"required,uuid"`
	Date            string  `json:"date" binding:"required"`
	CheckIn         *string `json:"checkIn"`
	CheckOut        *string `json:"checkOut"`
	ShiftScheduleID *string `json:"shiftScheduleId" binding:"omitempty,uuid"`
	Notes           *string `json:"notes"`
}

type ConfirmImageRequest struct {
	Key string `json:"key" binding:"required"`
}

func ToAttendanceDomain(req AttendanceRequest) (domain.Attendance, error) {
	empID, _ := uuid.Parse(req.EmployeeID)

	var shiftID *uuid.UUID
	if req.ShiftScheduleID != nil {
		id, _ := uuid.Parse(*req.ShiftScheduleID)
		shiftID = &id
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return domain.Attendance{}, fmt.Errorf("invalid date format: %w", err)
	}

	var checkIn *time.Time
	if req.CheckIn != nil {
		t, err := time.Parse("2006-01-02 15:04:05", *req.CheckIn)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04", *req.CheckIn)
			if err != nil {
				return domain.Attendance{}, fmt.Errorf("invalid checkIn format: %w", err)
			}
		}
		checkIn = &t
	}

	var checkOut *time.Time
	if req.CheckOut != nil {
		t, err := time.Parse("2006-01-02 15:04:05", *req.CheckOut)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04", *req.CheckOut)
			if err != nil {
				return domain.Attendance{}, fmt.Errorf("invalid checkOut format: %w", err)
			}
		}
		checkOut = &t
	}

	return domain.Attendance{
		EmployeeID:      empID,
		Date:            date,
		CheckIn:         checkIn,
		CheckOut:        checkOut,
		ShiftScheduleID: shiftID,
		Notes:           req.Notes,
	}, nil
}

func ToAttendanceUpdateResponse(a domain.Attendance) AttendanceResponse {
	var empCode, empName string
	if a.Employee != nil {
		empCode = a.Employee.Code
		empName = a.Employee.Name
	}

	dateStr := a.Date.Format("2006-01-02")

	var checkInStr string
	if a.CheckIn != nil {
		checkInStr = a.CheckIn.Format("15:04")
	}

	var checkOutStr *string
	var totalWorkMinsStr *string

	if a.CheckOut != nil {
		co := a.CheckOut.Format("15:04")
		checkOutStr = &co

		if a.CheckIn != nil {
			duration := a.CheckOut.Sub(*a.CheckIn)
			strMins := fmt.Sprintf("%d", int(duration.Minutes()))
			totalWorkMinsStr = &strMins
		}
	}

	var shiftName string
	if a.ShiftSchedule != nil {
		shiftName = a.ShiftSchedule.Name
	}

	var delAt *string
	if a.DeletedAt.Valid {
		t := a.DeletedAt.Time.Format("2006-01-02 15:04:05")
		delAt = &t
	}

	return AttendanceResponse{
		ID:                a.ID,
		EmployeeCode:      empCode,
		EmployeeName:      empName,
		Date:              dateStr,
		CheckIn:           checkInStr,
		CheckOut:          checkOutStr,
		TotalWorkMinutes:  totalWorkMinsStr,
		ShiftScheduleName: shiftName,
		Notes:             a.Notes,
		LateMinutes:       a.LateMinutes,
		DeductionAmount:   decimal.NewFromFloat(a.DeductionAmount),
		DeletedAt:         delAt,
		PhotoKey:          a.PhotoKey,
	}
}

func ToCreateAttedanceResponse(a domain.Attendance, uploadUrl string) AttendanceResponse {
	resp := ToAttendanceUpdateResponse(a)
	resp.UploadURL = uploadUrl
	return resp
}

func ToAttendanceResponses(as []domain.Attendance) []AttendanceResponse {
	res := make([]AttendanceResponse, 0, len(as))

	for _, a := range as {
		res = append(res, ToAttendanceUpdateResponse(a))
	}

	return res
}

func ToAttendanceResponsePagination(a []AttendanceResponse, p filter.PaginationWithInputFilter, totalRows int64) AttendanceResponsePagination {
	return AttendanceResponsePagination{
		Data: a,
		Meta: p.ToMeta(totalRows),
	}
}
