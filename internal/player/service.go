package player

import (
	"Go-lab/internal/utils"
	"context"
	"database/sql"
)

type Service struct {
	db      *utils.DbUtils
	repo    *Repo
	api     *API
	ctx     context.Context
}

func NewService(dbUtils *utils.DbUtils, repo *Repo, api *API) *Service {
	ctx := context.Background()
	service := &Service{
		db:     dbUtils,
		repo:   repo,
		api:    api,
		ctx:    ctx,
	}
	return service
}


func (s *Service) FindAll(ctx context.Context) ([]Player, error) {
	var res []Player
	
	err := s.db.WithTransaction(func(tx *sql.Tx) error {
		var err error
		res, err = s.repo.FindAll(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return res, nil
}

func (s *Service) FindById(ctx context.Context, id int) (*Player, error) {
	var playerEntity *Player

	err := s.db.WithTransaction(func(tx *sql.Tx) error {
		p, err := s.repo.FindById(ctx, id)
		if err != nil {
			return err
		}
		playerEntity = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return playerEntity, nil
}

func (s *Service) FindByResourceId(ctx context.Context, resourceId string) (*Player, error) {
	var playerEntity *Player

	err := s.db.WithTransaction(func(tx *sql.Tx) error {
		p, err := s.repo.FindByResourceId(ctx, resourceId)
		if err != nil {
			return err
		}
		playerEntity = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return playerEntity, nil
}

func (s *Service) Checkin(ctx context.Context, id int) (*Player, error) {
	var playerEntity *Player

	err := s.db.WithTransaction(func(tx *sql.Tx) error {
		p, err := s.repo.Checkin(ctx, id)
		if err != nil {
			return err
		}
		playerEntity = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return playerEntity, nil
}