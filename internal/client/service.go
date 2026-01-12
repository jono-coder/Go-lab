package client

import (
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/dbutils"
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Service struct {
	db      *dbutils.DbUtils
	repo    *Repo
	api     *API
	running atomic.Bool
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewService(dbUtils *dbutils.DbUtils, repo *Repo, api *API) *Service {
	if err := validate.Required("dbUtils", dbUtils); err != nil {
		panic(err)
	}
	if err := validate.Required("repo", repo); err != nil {
		panic(err)
	}
	if err := validate.Required("api", api); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	service := &Service{
		db:     dbUtils,
		repo:   repo,
		api:    api,
		ctx:    ctx,
		cancel: cancel,
	}
	return service
}

func (s *Service) FindById(ctx context.Context, id uint) (*Client, error) {
	if err := validate.Required("ctx", ctx); err != nil {
		return nil, err
	}

	var client *Client

	err := s.db.WithTransaction(ctx, func(*sql.Tx) error {
		//userId, _ := session.UserIDFromContext(ctx)
		//log.Println("userId:", userId)

		c, err := s.repo.FindById(ctx, id)
		if err != nil {
			return err
		}
		client = c
		return nil
	})

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *Service) FindAll(ctx context.Context) ([]Client, error) {
	if err := validate.Required("ctx", ctx); err != nil {
		return nil, err
	}

	var clients []Client

	if err := s.db.WithTransaction(ctx, func(*sql.Tx) error {
		var err error
		if clients, err = s.repo.FindAll(ctx, utils.NewPaging()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return clients, nil
}

func (s *Service) DoBusinessStuff(ctx context.Context) error {
	if err := validate.Required("ctx", ctx); err != nil {
		return err
	}

	err := s.db.WithTransaction(ctx, func(*sql.Tx) error {
		var err error
		clients, err := s.repo.FindAll(ctx, utils.NewPaging())
		if err != nil {
			return err
		}
		log.Printf("we found '%d' clients", len(clients))

		clientEntity, err := s.repo.FindById(ctx, 1)
		if err != nil {
			return err
		}
		log.Printf("we also found our clientEntity: a/c no = '%v'", clientEntity.AccountNo)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetById(id uint) (*Client, error) {
	res, _, err := s.api.GetById(id)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (s *Service) GetAll() ([]Client, error) {
	res, _, err := s.api.GetAll()
	return res, err
}

func (s *Service) Name() string {
	return "ClientService"
}

func (s *Service) Start() {
	if !s.running.CompareAndSwap(false, true) {
		return
	}

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				log.Println("Client Scheduler running...")
				time.Sleep(5 * time.Second)
				log.Println("Client Scheduler ran.")
			}
		}
	}()
}

func (s *Service) Stop() {
	go func() {
		log.Println("Stopping Client Scheduler...")

		if !s.running.CompareAndSwap(true, false) {
			return
		}

		s.cancel()
		s.wg.Wait()

		log.Println("Client Scheduler stopped.")
	}()
}

func (s *Service) IsRunning() bool {
	return s.running.Load()
}
