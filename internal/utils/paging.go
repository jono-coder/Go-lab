package utils

const (
	LimitDefault  = 20
	OffsetDefault = 0
)

type Paging struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PagingOption func(*Paging)

func WithLimit(limit int) PagingOption {
	return func(p *Paging) {
		if limit > 0 {
			p.Limit = limit
		}
	}
}

func WithOffset(offset int) PagingOption {
	return func(p *Paging) {
		if offset >= 0 {
			p.Offset = offset
		}
	}
}

func NewPaging(opts ...PagingOption) Paging {
	p := Paging{
		Limit:  LimitDefault,
		Offset: OffsetDefault,
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}
