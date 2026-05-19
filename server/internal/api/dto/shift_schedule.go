package dto

import (
	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ShiftScheduleResponse struct {
	ID                  uuid.UUID       `json:"id"`
	Name                string          `json:"name"`
	StartTime           string          `json:"startTime"`
	EndTime             string          `json:"endTime"`
	ToleranceMinutes    int             `json:"toleranceMinutes"`
	LateIntervalMinutes int             `json:"lateIntervalMinutes"`
	LateDeductionAmount decimal.Decimal `json:"lateDeductionAmount"`
	DeletedAt           *string         `json:"deletedAt,omitempty"`
}

type ShiftScheduleResponsePagination struct {
	Data []ShiftScheduleResponse `json:"data"`
	Meta filter.Meta             `json:"meta"`
}

type ShiftScheduleRequest struct {
	Name                string  `json:"name" binding:"required"`
	StartTime           string  `json:"startTime" binding:"required"`
	EndTime             string  `json:"endTime" binding:"required"`
	ToleranceMinutes    int     `json:"toleranceMinutes" binding:"min=0"`
	LateIntervalMinutes int     `json:"lateIntervalMinutes" binding:"min=1"`
	LateDeductionAmount float64 `json:"lateDeductionAmount" binding:"min=0"`
}

func ToShiftScheduleResponse(s domain.ShiftSchedule) ShiftScheduleResponse {
	res := ShiftScheduleResponse{
		ID:                  s.ID,
		Name:                s.Name,
		StartTime:           s.StartTime,
		EndTime:             s.EndTime,
		ToleranceMinutes:    s.ToleranceMinutes,
		LateIntervalMinutes: s.LateIntervalMinutes,
		LateDeductionAmount: decimal.NewFromFloat(s.LateDeductionAmount),
	}

	if s.DeletedAt.Valid {
		t := s.DeletedAt.Time.Format("2006-01-02 15:04:05")
		res.DeletedAt = &t
	}

	return res
}

func ToShiftScheduleResponses(ss []domain.ShiftSchedule) []ShiftScheduleResponse {
	res := make([]ShiftScheduleResponse, 0, len(ss))
	for _, s := range ss {
		res = append(res, ToShiftScheduleResponse(s))
	}
	return res
}

func ToShiftScheduleResponsePagination(ss []ShiftScheduleResponse, p filter.PaginationWithInputFilter, totalRows int64) ShiftScheduleResponsePagination {
	return ShiftScheduleResponsePagination{
		Data: ss,
		Meta: p.ToMeta(totalRows),
	}
}

func ToShiftScheduleDomain(req ShiftScheduleRequest) domain.ShiftSchedule {
	return domain.ShiftSchedule{
		Name:                req.Name,
		StartTime:           req.StartTime,
		EndTime:             req.EndTime,
		ToleranceMinutes:    req.ToleranceMinutes,
		LateIntervalMinutes: req.LateIntervalMinutes,
		LateDeductionAmount: req.LateDeductionAmount,
	}
}
