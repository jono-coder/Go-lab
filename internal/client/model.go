package client

import (
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/validate"
	"time"
)

type Client struct {
	Id          *uint
	AccountNo   string `validate:"required,notblank,alphanumunicode"`
	AccountName string `validate:"required,notblank,min=3,max=50"`
	CreatedAt   *time.Time
}

func NewClient(accountNo, accountName string) (*Client, error) {
	c := Client{
		AccountNo:   accountNo,
		AccountName: accountName,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Client) Validate() error {
	return validate.Get().Struct(c)
}

func (c *Client) String() string {
	return utils.ToString(c)
}

// AvgSize used with caching
func AvgSize() int64 {
	return 8 + 10 + 30 + 8 // add all fields len
}
