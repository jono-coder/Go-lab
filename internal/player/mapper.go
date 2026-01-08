package player

import "time"

func ToDTO(p *Player) DTO {
	var lc *time.Time
	if p.LastCheckin.Valid {
		lc = &p.LastCheckin.Time
	}
	return DTO{
		Id:          p.Id,
		ResourceId:  p.ResourceId,
		Name:        p.Name,
		Description: p.Description,
		LastCheckin: lc,
		CreatedAt:   p.CreatedAt,
	}
}

func ToDTOs(players []Player) []DTO {
	res := make([]DTO, 0, len(players))
	for i, player := range players {
		res[i] = ToDTO(&player)
	}
	return res
}
