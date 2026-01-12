package player

func ToDTO(p *Player) DTO {
	return DTO{
		Id:          p.Id,
		ResourceId:  p.ResourceId,
		Name:        p.Name,
		Description: *p.Description,
		LastCheckin: p.LastCheckin,
		CreatedAt:   p.CreatedAt,
	}
}

func ToDTOs(players []Player) []DTO {
	res := make([]DTO, len(players))
	for i, player := range players {
		res[i] = ToDTO(&player)
	}
	return res
}
