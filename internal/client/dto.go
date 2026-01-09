package client

import "time"

type DTO struct {
	Id          *int       `json:"id"`
	AccountNo   string     `json:"account_no"`
	AccountName string     `json:"account_name"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}
