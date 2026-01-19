package player

import (
	"Go-lab/internal/utils"
	"time"
)

/*
	notes:
	Removed omitempty for LastCheckin and CreatedAt to ensure a stable frontend contract. TBD
*/

type DTO struct {
	Id          *uint      `json:"id"`
	ResourceId  string     `json:"resource_id" `
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	LastCheckin *time.Time `json:"last_checkin"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   *uint      `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *uint      `json:"updated_by"`
}

func (d *DTO) String() string {
	return utils.ToString(d)
}
