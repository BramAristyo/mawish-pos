package database

import (
	"strings"
	"time"

	"github.com/BramAristyo/mawish-pos/server/pkg/filter"
	"gorm.io/gorm"
)

type PreloadEntity struct {
	Entity string
}

func BuildQuery(
	db *gorm.DB,
	req filter.DynamicFilter,
	searchable []string,
	allowedFields map[string]string,
) *gorm.DB {
	db = applySearch(db, req.Search, searchable, allowedFields)
	db = applySort(db, req.Sort, allowedFields)
	db = applyFilter(db, req.Filter, allowedFields)

	return db.Debug()
}

func applySearch(
	db *gorm.DB,
	search string,
	searchable []string,
	allowedFields map[string]string,
) *gorm.DB {
	if search == "" {
		return db
	}

	conditions := make([]string, 0, len(searchable))
	args := make([]any, 0, len(searchable))

	for _, field := range searchable {
		if col, ok := allowedFields[field]; ok {
			conditions = append(conditions, col+" ILIKE ?")
			args = append(args, "%"+search+"%")
		}
	}

	if len(conditions) > 0 {
		db = db.Where("("+strings.Join(conditions, " OR ")+")", args...)
	}

	return db
}

func applySort(
	db *gorm.DB,
	sorts []filter.Sort,
	allowedFields map[string]string,
) *gorm.DB {
	for _, s := range sorts {
		col, ok := allowedFields[s.Column]
		if !ok {
			continue
		}

		order := strings.ToLower(s.Order)
		if order != "asc" && order != "desc" {
			continue
		}

		db = db.Order(col + " " + order)
	}

	return db
}

func applyFilter(
	db *gorm.DB,
	filters map[string]filter.Filter,
	allowedFields map[string]string,
) *gorm.DB {

	for field, f := range filters {

		col, ok := allowedFields[field]
		if !ok {
			continue
		}

		if f.FilterType != filter.DataTypeDate {
			continue
		}

		if f.Type == filter.OpInRange && f.From != "" && f.To != "" {

			fromDate, err := time.Parse("2006-01-02", f.From)
			if err != nil {
				continue
			}

			toDate, err := time.Parse("2006-01-02", f.To)
			if err != nil {
				continue
			}

			toDate = toDate.Add(24*time.Hour - time.Nanosecond)

			db = db.Where(col+" BETWEEN ? AND ?", fromDate, toDate)
		}

		if f.Type == filter.OpGreaterThanOrEqual && f.From != "" {

			fromDate, err := time.Parse("2006-01-02", f.From)
			if err != nil {
				continue
			}

			db = db.Where(col+" >= ?", fromDate)
		}

		if f.Type == filter.OpLessThanOrEqual && f.From != "" {

			toDate, err := time.Parse("2006-01-02", f.From)
			if err != nil {
				continue
			}

			toDate = toDate.Add(24*time.Hour - time.Nanosecond)

			db = db.Where(col+" <= ?", toDate)
		}

		if f.Type == filter.OpGreaterThan && f.From != "" {

			fromDate, err := time.Parse("2006-01-02", f.From)
			if err != nil {
				continue
			}

			db = db.Where(col+" > ?", fromDate)
		}

		if f.Type == filter.OpLessThan && f.From != "" {

			toDate, err := time.Parse("2006-01-02", f.From)
			if err != nil {
				continue
			}

			db = db.Where(col+" < ?", toDate)
		}
	}

	return db
}
