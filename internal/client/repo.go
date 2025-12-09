package client

import (
	"Go-lab/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type Repo struct {
	db    *utils.DbUtils
	cache *ristretto.Cache[int, Client]
}

func NewRepo(dbUtils *utils.DbUtils) *Repo {
	cache, err := ristretto.NewCache(&ristretto.Config[int, Client]{
		NumCounters: 10_000,  // number of keys to track frequency of (10M).
		MaxCost:     1 << 25, // maximum cost of cache (32MB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Repo{
		db:    dbUtils,
		cache: cache,
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func (r *Repo) FindById(ctx context.Context, id int) (*Client, error) {
	if client, found := r.cache.Get(id); found {
		return &client, nil
	}

	var res Client

	err := r.db.DB.QueryRowContext(ctx,
		`SELECT id, account_no, account_name, created_at
             	   FROM client_entity
             	   WHERE id = ?`, id,
	).Scan(&res.Id, &res.AccountNo, &res.AccountName, &res.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("find client %d: %w", id, err)
	}

	if wasSet := r.cache.SetWithTTL(id, res, AvgSize(), 5*time.Minute); !wasSet {
		log.Println("Was not added to the cache", wasSet)
	}
	r.cache.Wait()

	return &res, nil
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func (r *Repo) FindAll(ctx context.Context) ([]Client, error) {
	var res []Client

	rows, err := r.db.DB.QueryContext(ctx,
		`SELECT id, account_no, account_name, created_at
             	   FROM client_entity
             	   ORDER BY account_no`,
	)
	if err != nil {
		return nil, fmt.Errorf("find all clients: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.Id, &client.AccountNo, &client.AccountName, &client.CreatedAt); err != nil {
			return nil, fmt.Errorf("find all clients: %w", err)
		}
		res = append(res, client)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (r *Repo) Close() {
	r.cache.Close()
}
