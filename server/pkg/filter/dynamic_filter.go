package filter

import (
	"time"

	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
)

type FilterOperator string

const (
	OpContains           FilterOperator = "contains"
	OpNotContains        FilterOperator = "notContains"
	OpEquals             FilterOperator = "equals"
	OpNotEqual           FilterOperator = "notEqual"
	OpStartsWith         FilterOperator = "startsWith"
	OpEndsWith           FilterOperator = "endsWith"
	OpLessThan           FilterOperator = "lessThan"
	OpLessThanOrEqual    FilterOperator = "lessThanOrEqual"
	OpGreaterThan        FilterOperator = "greaterThan"
	OpGreaterThanOrEqual FilterOperator = "greaterThanOrEqual"
	OpInRange            FilterOperator = "inRange"
)

type FilterDataType string

const (
	DataTypeText   FilterDataType = "text"
	DataTypeNumber FilterDataType = "number"
	DataTypeDate   FilterDataType = "date"
)

type Sort struct {
	Column string `json:"column" form:"column"`
	Order  string `json:"orderBy" form:"orderBy"`
}

type Filter struct {
	Type       FilterOperator `json:"type" form:"type"`
	From       string         `json:"from" form:"from"`
	To         string         `json:"to" form:"to"`
	FilterType FilterDataType `json:"filterType" form:"filterType"`
}

type DynamicFilter struct {
	Sort   []Sort            `json:"sort" form:"sort"`
	Search string            `json:"search" form:"search"`
	Filter map[string]Filter `json:"filter" form:"filter"`
}

func (df *DynamicFilter) WithDefaultSort() {
	if len(df.Sort) == 0 {
		df.Sort = []Sort{{"created_at", "desc"}}
	}

}

func (df *DynamicFilter) WithDefaultDateRange() {
	now := time.Now()

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	from := startOfMonth.Format("2006-01-02")
	to := now.Format("2006-01-02")

	if df.Filter == nil {
		df.Filter = make(map[string]Filter)
	}

	if _, exists := df.Filter["created_at"]; !exists {

		df.Filter["created_at"] = Filter{
			Type:       OpInRange,
			From:       from,
			To:         to,
			FilterType: DataTypeDate,
		}
	}
}

func (df *DynamicFilter) ValidateHasDateRange() error {
	if len(df.Filter) == 0 {
		return usecase_errors.DateFilterRequired
	}

	hasValidDateRange := false

	for _, f := range df.Filter {
		if f.FilterType == DataTypeDate && f.Type == OpInRange {
			if f.From != "" && f.To != "" {
				hasValidDateRange = true
				break
			}
		}
	}

	if !hasValidDateRange {
		return usecase_errors.DateFilterRequired
	}

	return nil
}
