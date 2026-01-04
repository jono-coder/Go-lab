package client

func ToDTO(c *Client) DTO {
	return DTO{
		Id:          c.Id,
		AccountNo:   c.AccountNo,
		AccountName: c.AccountName,
		CreatedAt:   c.CreatedAt,
	}
}

func ToDTOs(clients []Client) []DTO {
	res := make([]DTO, len(clients))
	for i, client := range clients {
		res[i] = ToDTO(&client)
	}
	return res
}
