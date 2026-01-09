package client

import (
	"Go-lab/config"
	"Go-lab/internal/security"
	"context"
	"log"
	"sync"
)

var (
	clientApi *API
	once      sync.Once
)

/*
func TestGetById_NotFound(t *testing.T) {
	beforeEach()
	defer afterEach()

	entity, code, err := clientApi.GetById(-1)

	req := require.New(t)
	req.Nil(entity)
	req.NotNil(err)
	req.Equal(http.StatusNotFound, code)
}


func TestGetById_Found(t *testing.T) {
	beforeEach()
	defer afterEach()

	entity, code, err := clientApi.GetById(1)
	if err != nil {
		t.Fail()
	} else {
		req := require.New(t)
		req.NotNil(entity)
		req.Equal(http.StatusOK, code)
	}
}

func TestGetAll(t *testing.T) {
	beforeEach()
	defer afterEach()

	entities, code, err := clientApi.GetAll()
	if err != nil {
		t.Fail()
	} else {
		req := require.New(t)
		req.NotNil(entities)
		req.NotEmpty(entities)
		req.Equal(http.StatusOK, code)
	}
}*/

func beforeEach() {
	var baseUrl string
	once.Do(func() {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("failed to load cfg: %v", err)
		}
		baseUrl = cfg.App.BaseUrl
	})

	var err error
	clientApi, err = NewAPI(security.NewOAuthConfig(context.Background(), baseUrl))
	if err != nil {
		log.Fatalf("failed: %v", err)
	}
}

func afterEach() {
}
