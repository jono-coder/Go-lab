package player

func ToDTO(p *Player) DTO {
	return DTO{
		Id:          p.Id,
		ResourceId:  p.ResourceId,
		Name:        p.Name,
		Description: p.Description,
		LastCheckin: p.LastCheckin,
		CreatedAt:   p.CreatedAt,
	}
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

func ToDTOs(players []Player) []DTO {
	res := make([]DTO, len(players))
	for i := range players {
		res[i] = ToDTO(&players[i])
	}
	return res
}
