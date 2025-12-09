package client

import (
	"Go-lab/internal/utils"
	"time"
)

type Client struct {
	Id          int       `json:"id"`
	AccountNo   string    `json:"account_no"`
	AccountName string    `json:"account_name"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *Client) NewClient() *Client {
	return &Client{}
}

func (c *Client) String() string {
	return utils.ToString(c)
}

// AvgSize used with caching
func AvgSize() int64 {
	return 8 + 10 + 30 + 8 // add all fields len
}
