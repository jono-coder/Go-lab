package player

import (
	"Go-lab/internal/security"
	"Go-lab/internal/utils/validate"

	"fmt"
)

type API struct {
	config *security.OAuthConfig
}

func NewAPI(config *security.OAuthConfig) (*API, error) {
	if err := validate.Get().Var(config, "required"); err != nil {
		return nil, err
	}

	return &API{config: config}, nil
}

func (c *API) GetAll() ([]Player, int, error) {
	var res []Player

	resp, err := c.config.Client.R().
		SetResult(&res).
		Get("player")
	if err != nil {
		return res, 0, err
	}
	if resp.IsError() {
		return res, resp.StatusCode(), fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}

	return res, resp.StatusCode(), nil
}

func (c *API) GetById(id int) (*Player, int, error) {
	var res Player

	resp, err := c.config.Client.R().
		SetResult(&res).
		Get(fmt.Sprintf("player/%d", id))
	if err != nil {
		return nil, 0, err
	}
	if resp.IsError() {
		return nil, resp.StatusCode(), fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}

	return &res, resp.StatusCode(), nil
}

func (c *API) GetByResourceId(resourceId string) (*Player, int, error) {
	var res Player

	resp, err := c.config.Client.R().
		SetResult(&res).
		Get(fmt.Sprintf("player/resource/%s", resourceId))
	if err != nil {
		return nil, 0, err
	}
	if resp.IsError() {
		return nil, resp.StatusCode(), fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}

	return &res, resp.StatusCode(), nil
}

func (c *API) Checkin(id int) (*Player, int, error) {
	var res Player

	resp, err := c.config.Client.R().
		SetResult(&res).
		Get(fmt.Sprintf("player/checkin/%d", id))
	if err != nil {
		return nil, 0, err
	}
	if resp.IsError() {
		return nil, resp.StatusCode(), fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}

	return &res, resp.StatusCode(), nil
}
