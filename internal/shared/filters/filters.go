// File: internal/data/shared/filters/filters.go

package filters

import (
	"strings"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
)

// Filters is a struct that holds filter parameters for querying data.
type Filters struct {
	Page         int      // Current page number
	PageSize     int      // Number of records per page
	Sort         string   // Sort parameter
	SortSafelist []string // List of permitted sort values
}

// MetaData holds pagination metadata.
type MetaData struct {
	CurrentPage  int `json:"current_page,omitempty"`  // Current page number
	PageSize     int `json:"page_size,omitempty"`     // Number of records per page
	FirstPage    int `json:"first_page,omitempty"`    // First page number
	LastPage     int `json:"last_page,omitempty"`     // Last page number
	TotalRecords int `json:"total_records,omitempty"` // Total number of records
}

// for empty MetaData
var EmptyMetaData = MetaData{}

// ValidateFilters checks the validity of the filter parameters.
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")                        // Page must be greater than 0
	v.Check(f.Page <= 500, "page", "must be a maximum of 500")                      // Page must be at most 500
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")               // PageSize must be greater than 0
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")             // PageSize must be at most 100
	v.Check(v.IsPermitted(f.Sort, f.SortSafelist...), "sort", "invalid sort value") // Sort must be in the safelist
}

// limit returns the limit for SQL queries based on the PageSize.
func (f Filters) Limit() int {
	return f.PageSize
}

// offset returns the offset for SQL queries based on the Page and PageSize.
func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// sortColumn returns the column name to sort by, trimsming any leading '-' for descending order.
func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-") // Remove leading '-' if present
		}
	}
	panic("unsafe sort parameter: " + f.Sort) // Panic if the sort parameter is not in the safelist
}

// sortDirection returns the sort direction (asc/desc) based on the Sort field.
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC" // Descending order if Sort starts with '-'
	}
	return "ASC" // Ascending order otherwise
}

// calculateMetadata computes pagination metadata based on total records.
func CalculateMetaData(totalRecords, page, pageSize int) MetaData {
	if totalRecords == 0 {
		return MetaData{} // Return empty metadata if there are no records
	}

	return MetaData{
		CurrentPage:  page,                                     // Current page number
		PageSize:     pageSize,                                 // Page size
		FirstPage:    1,                                        // First page is always 1
		LastPage:     (totalRecords + pageSize - 1) / pageSize, // Ceiling division
		TotalRecords: totalRecords,                             // Total number of records
	}
}
