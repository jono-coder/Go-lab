package player

import (
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/validate"
	"time"
)

type Player struct {
	Id          *int
	ResourceId  string
	Name        string
	Description *string
	LastCheckin *time.Time
	CreatedAt   *time.Time
}

func NewPlayer(resourceId, name string, description *string) (*Player, error) {
	if err := validate.NotBlank("resourceId", resourceId); err != nil {
		return nil, err
	}
	if err := validate.NotBlank("name", name); err != nil {
		return nil, err
	}

	return &Player{
		ResourceId:  resourceId,
		Name:        name,
		Description: description,
	}, nil
}

func (p *Player) String() string {
	return utils.ToString(p)
}
