package client

import (
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/dbutils"
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type Repo struct {
	db    *dbutils.DbUtils
	cache *ristretto.Cache[uint, Client]
}

func NewRepo(dbUtils *dbutils.DbUtils) *Repo {
	if err := validate.Required("dbUtils", dbUtils); err != nil {
		panic(err)
	}

	cache, err := ristretto.NewCache(&ristretto.Config[uint, Client]{
		NumCounters: 10_000,   // number of keys to track frequency of (10M).
		MaxCost:     32 << 20, // maximum cost of cache (32MB).
		BufferItems: 64,       // number of keys per Get buffer.
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
func (r *Repo) FindById(ctx context.Context, id uint) (*Client, error) {
	if client, found := r.cache.Get(id); found {
		return &client, nil
	}

	if err := validate.Required("ctx", ctx); err != nil {
		return nil, err
	}

	var (
		accountNo   string
		accountName string
		createdAt   time.Time
	)

	if err := r.db.DB.QueryRowContext(ctx,
		`SELECT
    				account_no,
    				account_name,
    				created_at
             	FROM
             		client_entity
             	WHERE
             	    id = ?`, id,
	).Scan(&accountNo, &accountName, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("find client %d: %w", id, err)
	}

	c, err := NewClient(accountNo, accountName)
	if err != nil {
		return nil, err
	}
	c.Id = &id
	c.CreatedAt = &createdAt

	if wasSet := r.cache.SetWithTTL(id, *c, AvgSize(), 5*time.Minute); !wasSet {
		log.Println("Was not added to the cache", wasSet)
	}
	r.cache.Wait()

	return c, nil
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func (r *Repo) FindAll(ctx context.Context, paging utils.Paging) ([]Client, error) {
	if err := validate.Required("ctx", ctx); err != nil {
		return nil, err
	}
	if err := validate.Required("paging", paging); err != nil {
		return nil, err
	}

	var clients []Client

	rows, err := r.db.DB.QueryContext(ctx,
		`SELECT
    				id,
    				account_no,
    				account_name,
    				created_at
             	FROM
             		client_entity
             	ORDER BY
             	    account_no ASC
                LIMIT ? OFFSET ?`, paging.Limit, paging.Offset,
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
		clients = append(clients, client)
	}

	if err = rows.Err(); err != nil {
		return clients, err
	}

	return clients, nil
}

func (r *Repo) Close() {
	r.cache.Close()
}
