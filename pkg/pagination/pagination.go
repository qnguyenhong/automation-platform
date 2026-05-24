package pagination

import "math"

const DefaultPerPage = 20
const MaxPerPage = 100

// Params holds pagination request parameters.
type Params struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// Offset returns the SQL OFFSET value.
func (p Params) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// Sanitize normalizes page and per_page values.
func (p Params) Sanitize() Params {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = DefaultPerPage
	}
	if p.PerPage > MaxPerPage {
		p.PerPage = MaxPerPage
	}
	return p
}

// Meta holds pagination response metadata.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewMeta creates pagination metadata.
func NewMeta(page, perPage int, total int64) Meta {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPages < 1 {
		totalPages = 1
	}
	return Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
