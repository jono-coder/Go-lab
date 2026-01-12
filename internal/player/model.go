package player

import (
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/validate"
	"time"
)

type Player struct {
	Id          *uint
	ResourceId  string  `validate:"required,notblank,alphanum,max=100"`
	Name        string  `validate:"required,notblank,min=1,max=50"`
	Description *string `validate:"min=1,max=50"`
	LastCheckin *time.Time
	CreatedAt   *time.Time
}

func NewPlayer(resourceId, name string, description *string) (*Player, error) {
	p := Player{
		ResourceId:  resourceId,
		Name:        name,
		Description: description,
	}

	if err := validate.Get().Struct(p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *Player) Validate() error {
	return validate.Get().Struct(p)
}

func (p *Player) String() string {
	return utils.ToString(p)
}
