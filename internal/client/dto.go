package client

import (
	"Go-lab/internal/utils"
	"time"
)

type DTO struct {
	Id          *uint      `json:"id"`
	AccountNo   string     `json:"account_no"`
	AccountName string     `json:"account_name"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

func (d *DTO) String() string {
	return utils.ToString(d)
}
