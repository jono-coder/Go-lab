package player

import (
	"database/sql"
	"time"
)

type Player struct {
	Id int
	ResourceId string
	Name string
	Description string
	LastCheckin sql.NullTime
	CreatedAt time.Time

}

func NewPlayer(resourceId string, name string, description string) *Player {
	return &Player{
		ResourceId: resourceId,
		Name: name,
		Description: description,
	}
}
