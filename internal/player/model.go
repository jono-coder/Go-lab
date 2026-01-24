package player

import (
	"Go-lab/internal/audit"
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/validate"
	"time"
)

type Player struct {
	audit.Auditable
	Id          *uint      `db:"id"`
	ResourceId  string     `db:"resource_id" validate:"required,notblank,max=100"`
	Name        string     `db:"name" validate:"required,notblank,max=50"`
	Description *string    `validate:"min=1,max=50"`
	LastCheckin *time.Time `db:"last_checkin"`
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
