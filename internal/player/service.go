package player

import (
	"Go-lab/internal/utils/dbutils"
	"Go-lab/internal/utils/paging"
	"Go-lab/internal/utils/validate"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	db  *dbutils.DbUtils
	api *API
	ctx context.Context
}

func NewService(dbUtils *dbutils.DbUtils, api *API) *Service {
	ctx := context.Background()
	service := &Service{
		db:  dbUtils,
		api: api,
		ctx: ctx,
	}
	return service
}

func (s *Service) Create(ctx context.Context, player *Player) (*uint, error) {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return nil, err
	}
	if err := validate.Get().Var(player, "required"); err != nil {
		return nil, err
	}

	var id *uint

	err := s.db.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		repo, err := s.createPlayerRepo(tx)
		if err != nil {
			return err
		}

		id, err = repo.Create(ctx, player)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *Service) FindAll(ctx context.Context, paging paging.Paging) ([]Player, error) {
	var players []Player

	err := s.db.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		var err error
		repo, err := s.createPlayerRepo(tx)
		if err != nil {
			return err
		}

		if players, err = repo.FindAll(ctx, paging); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return players, nil
}

func (s *Service) FindById(ctx context.Context, id uint) (*Player, error) {
	var playerEntity *Player

	err := s.db.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		repo, err := s.createPlayerRepo(tx)
		if err != nil {
			return err
		}

		p, err := repo.FindById(ctx, id)
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
	var player *Player

	err := s.db.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		repo, err := s.createPlayerRepo(tx)
		if err != nil {
			return err
		}

		p, err := repo.FindByResourceId(ctx, resourceId)
		if err != nil {
			return err
		}
		player = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return player, nil
}

func (s *Service) Checkin(ctx context.Context, id uint, updatedAt *time.Time) (*Player, error) {
	var player *Player

	err := s.db.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		repo, err := s.createPlayerRepo(tx)
		if err != nil {
			return err
		}

		p, err := repo.Checkin(ctx, id, updatedAt)
		if err != nil {
			return err
		}
		player = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return player, nil
}

func (s *Service) Update(ctx context.Context, dto *UpdateDto) (*Player, error) {
	var player *Player

	err := s.db.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		repo, err := s.createPlayerRepo(tx)
		if err != nil {
			return err
		}

		p, err := repo.Update(ctx, dto)
		if err != nil {
			return err
		}
		player = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return player, nil
}

func (s *Service) createPlayerRepo(tx *sqlx.Tx) (*Repo, error) {
	repo, err := NewRepo(tx)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
