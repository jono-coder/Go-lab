package client

import (
	"Go-lab/internal/security"
	"Go-lab/internal/utils/validate"

	"fmt"
)

type API struct {
	config *security.OAuthConfig
}

func NewAPI(config *security.OAuthConfig) (*API, error) {
	err := validate.Required("config", config)
	if err != nil {
		return nil, err
	}

	return &API{config: config}, nil
}

func (c *API) GetAll() ([]Client, int, error) {
	var res []Client

	resp, err := c.config.Client.R().
		SetResult(&res).
		Get("client")
	if err != nil {
		return res, 0, err
	}
	if resp.IsError() {
		return res, resp.StatusCode(), fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}

	return res, resp.StatusCode(), nil
}

func (c *API) GetById(id int) (*Client, int, error) {
	var res Client

	resp, err := c.config.Client.R().
		SetResult(&res).
		Get(fmt.Sprintf("client/%d", id))
	if err != nil {
		return nil, 0, err
	}
	if resp.IsError() {
		return nil, resp.StatusCode(), fmt.Errorf("status %d: %s", resp.StatusCode(), resp.String())
	}

	return &res, resp.StatusCode(), nil
}
