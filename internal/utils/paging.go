package utils

const DefaultLimit uint = 20

type Paging struct {
	Page  uint `json:"page"` // 0-based
	Limit uint `json:"limit"`
}

func NewPaging(page uint, limit uint) Paging {
	if limit == 0 {
		limit = DefaultLimit
	}
	return Paging{Page: page, Limit: limit}
}

// Offset returns the starting index for the current page.
func (p Paging) Offset() uint {
	return p.Page * p.Limit
}
