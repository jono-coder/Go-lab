package contact

import (
	"Go-lab/internal/utils"
	"context"
	"database/sql"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Service struct {
	db      *utils.DbUtils
	repo    *Repo
	running atomic.Bool
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewService(dbUtils *utils.DbUtils, repo *Repo) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	service := &Service{
		db:     dbUtils,
		repo:   repo,
		ctx:    ctx,
		cancel: cancel,
	}
	return service
}

func (s *Service) FindById(ctx context.Context, id int) (*Contact, error) {
	var clientEntity *Contact

	err := s.db.WithTransaction(func(tx *sql.Tx) error {
		c, err := s.repo.FindById(ctx, id)
		if err != nil {
			return err
		}
		clientEntity = c
		return nil
	})

	if err != nil {
		return nil, err
	}

	return clientEntity, nil
}

func (s *Service) Name() string {
	return "ContactService"
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
				log.Println("Contact Scheduler running...")
				time.Sleep(3 * time.Second)
				log.Println("Contact Scheduler ran.")
			}
		}
	}()
}

func (s *Service) Stop() {
	go func() {
		log.Println("Stopping Contact Scheduler...")

		if !s.running.CompareAndSwap(true, false) {
			return
		}

		s.cancel()
		s.wg.Wait()

		log.Println("Contact Scheduler stopped.")
	}()
}

func (s *Service) IsRunning() bool {
	return s.running.Load()
}
