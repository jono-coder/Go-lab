package player

import (
	"Go-lab/internal/utils/validate"
	"fmt"
)

func ToDTO(p *Player) (*DTO, error) {
	if err := validate.Get().Var(p, "required"); err != nil {
		return nil, err
	}

	return &DTO{
		Id:          p.Id,
		ResourceId:  p.ResourceId,
		Name:        p.Name,
		Description: p.Description,
		LastCheckin: p.LastCheckin,
		CreatedAt:   p.CreatedAt,
		CreatedBy:   p.CreatedBy,
		UpdatedAt:   p.UpdatedAt,
		UpdatedBy:   p.UpdatedBy,
	}, nil
}

/*
	notes:

	Issues:
	for i, player := range players {
		res[i] = ToDTO(&player) // BUG
	}
	In Go, range copies each slice element into the variable player.
	Taking &player points to the same memory address on every iteration.
	All elements in res end up referencing the last loop iteration.
	This causes random fields to disappear or be overwritten, e.g., last_checkin appearing nil.

	What changed:
	&players[i] points to the actual slice element, not a copy.
	Each DTO now correctly references its player.
	Fixes “random missing field” bugs after refetches or updates.
*/

func ToDTOs(players []Player) ([]DTO, error) {
	res := make([]DTO, len(players))
	for i := range players {
		dto, err := ToDTO(&players[i])
		if err != nil {
			return nil, err
		}
		res[i] = *dto
	}
	return res, nil
}

func ToEntity(dto DTO) (*Player, error) {
	if err := validate.Get().Var(dto, "required"); err != nil {
		return nil, err
	}

	if dto.Id != nil {
		return nil, fmt.Errorf("id is not allowed to be set on creation")
	}

	player, err := NewPlayer(dto.ResourceId, dto.Name, dto.Description)
	if err != nil {
		return nil, err
	}

	player.LastCheckin = dto.LastCheckin

	if err = player.Validate(); err != nil {
		return nil, err
	}

	return player, nil
}
