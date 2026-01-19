package paging

import (
	"fmt"
	"strconv"
)

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

func ParsePage(pageStr string) (uint, error) {
	if pageStr == "" {
		return 0, nil
	}

	p, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, fmt.Errorf("invalid page number")
	}

	return uint(p), nil
}

func ParseLimit(limitStr string) (uint, error) {
	if limitStr == "" {
		return 0, nil
	}

	p, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, fmt.Errorf("invalid limit")
	}

	return uint(p), nil
}
