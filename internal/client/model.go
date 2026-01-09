package client

import (
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/validate"
	"time"
)

type Client struct {
	Id          *int
	AccountNo   string
	AccountName string
	CreatedAt   *time.Time
}

func NewClient(accountNo, accountName string) (*Client, error) {
	if err := validate.NotBlank("accountNo", accountNo); err != nil {
		return nil, err
	}

	if err := validate.NotBlank("accountName", accountName); err != nil {
		return nil, err
	}

	return &Client{
		AccountNo:   accountNo,
		AccountName: accountName,
	}, nil
}

func (c *Client) String() string {
	return utils.ToString(c)
}

// AvgSize used with caching
func AvgSize() int64 {
	return 8 + 10 + 30 + 8 // add all fields len
}
