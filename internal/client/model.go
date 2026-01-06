package client

import (
	"Go-lab/internal/utils"
	"time"
)

type Client struct {
	id          int
	AccountNo   string
	AccountName string
	CreatedAt   time.Time
}

func (c *Client) NewClient() *Client {
	return &Client{}
}

func (c *Client) Id() int {
	return c.id
}

func (c *Client) String() string {
	return utils.ToString(c)
}

// AvgSize used with caching
func AvgSize() int64 {
	return 8 + 10 + 30 + 8 // add all fields len
}