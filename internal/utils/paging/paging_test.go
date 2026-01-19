package paging

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPaging(t *testing.T) {
	var paging Paging

	req := require.New(t)

	paging = NewPaging(0, 0)
	req.NotNil(paging)
	req.Equal(DefaultLimit, paging.Limit)
	req.Equal(uint(0), paging.Offset())

	paging = NewPaging(0, 10)
	req.NotNil(paging)
	req.Equal(uint(10), paging.Limit)
	req.Equal(uint(0), paging.Offset())

	paging = NewPaging(1, 11)
	req.NotNil(paging)
	req.Equal(uint(11), paging.Limit)
	req.Equal(uint(11), paging.Offset())

	paging = NewPaging(1, 12)
	req.NotNil(paging)
	req.Equal(uint(12), paging.Limit)
	req.Equal(uint(12), paging.Offset())

	paging = NewPaging(2, 11)
	req.NotNil(paging)
	req.Equal(uint(11), paging.Limit)
	req.Equal(uint(22), paging.Offset())
}
