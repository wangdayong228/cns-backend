package utils

type Pagination struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}

func (p Pagination) CalcOffsetLimit() (offset int, limit int) {

	if p.Page <= 0 {
		p.Page = 1
	}

	switch {
	case p.PageSize > 100:
		p.PageSize = 100
	case p.PageSize <= 0:
		p.PageSize = 10
	}

	offset = (p.Page - 1) * p.PageSize
	limit = p.PageSize
	return
}
