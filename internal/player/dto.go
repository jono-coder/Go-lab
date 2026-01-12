package player

import (
	"Go-lab/internal/utils"
	"time"
)

type DTO struct {
	Id          *uint      `json:"id"`
	ResourceId  string     `json:"resource_id" `
	Name        string     `json:"name"`
	Description string     `json:"description"`
	LastCheckin *time.Time `json:"last_checkin,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

func (d *DTO) String() string {
	return utils.ToString(d)
}
