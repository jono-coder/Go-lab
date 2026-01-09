package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPagingDefault(t *testing.T) {
	paging := NewPaging()

	req := require.New(t)
	req.NotNil(paging)
	req.Equal(LimitDefault, paging.Limit)
	req.Equal(OffsetDefault, paging.Offset)
}

func TestNewPagingOk(t *testing.T) {
	limit := 100
	offset := 100
	paging := NewPaging(
		WithLimit(limit),
		WithOffset(offset))

	req := require.New(t)
	req.NotNil(paging)
	req.Equal(limit, paging.Limit)
	req.Equal(offset, paging.Offset)
}

func TestNewPagingIgnoredOptions(t *testing.T) {
	limit := -1
	offset := limit
	paging := NewPaging(
		WithLimit(limit),
		WithOffset(offset))

	req := require.New(t)
	req.NotNil(paging)
	req.Equal(LimitDefault, paging.Limit)
	req.Equal(OffsetDefault, paging.Offset)
}
