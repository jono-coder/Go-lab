package player

import "time"

type DTO struct {
	Id int `json:"id"`
	ResourceId string `json:"resource_id"`
	Name string `json:"name"`
	Description string `json:"description"`
	LastCheckin *time.Time `json:"last_checkin,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}