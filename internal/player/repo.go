package player

import (
	"Go-lab/internal/utils/dbutils"
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Repo struct {
	db *dbutils.DbUtils
}

func NewRepo(dbUtils *dbutils.DbUtils) *Repo {
	if err := validate.Required("dbUtils", dbUtils); err != nil {
		panic(err)
	}

	return &Repo{
		db: dbUtils,
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func (r *Repo) FindById(ctx context.Context, id uint) (*Player, error) {
	var res Player
	res.Id = &id

	err := r.db.DB.QueryRowContext(ctx,
		`SELECT
			resource_id,
			name,
			description,
			last_checkin,
			created_at
		FROM
			player_entity
        WHERE
			id = ?`,
		id,
	).Scan(&res.ResourceId, &res.Name, &res.Description, &res.LastCheckin, &res.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("find player %d: %w", id, err)
	}

	return &res, nil
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func (r *Repo) FindByResourceId(ctx context.Context, resourceId string) (*Player, error) {
	var (
		id          uint
		name        string
		description *string
		lastCheckin *time.Time
		createdAt   *time.Time
	)

	if err := r.db.DB.QueryRowContext(ctx,
		`SELECT
			id,
			name,
			description,
			last_checkin,
			created_at
		FROM
			player_entity
        WHERE
			resource_id = ?`,
		resourceId,
	).Scan(&id, &name, &description, &lastCheckin, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("find player %s: %w", resourceId, err)
	}

	res, err := NewPlayer(resourceId, name, description)
	if err != nil {
		return nil, err
	}
	res.Id = &id

	return res, nil
}

func (r *Repo) FindAll(ctx context.Context) ([]Player, error) {
	var res []Player

	rows, err := r.db.DB.QueryContext(ctx,
		`SELECT
			id,
			resource_id,
			name,
			description,
			last_checkin,
			created_at
		FROM
			player_entity
        ORDER BY
			name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("find all players: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var player Player
		if err := rows.Scan(&player.Id, &player.ResourceId, &player.Name, &player.Description, &player.LastCheckin, &player.CreatedAt); err != nil {
			return nil, fmt.Errorf("find all players: %w", err)
		}
		res = append(res, player)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (r *Repo) Checkin(ctx context.Context, id uint) (*Player, error) {
	res, err := r.db.DB.ExecContext(ctx,
		`UPDATE
			player_entity
		SET
			last_checkin = CURRENT_TIMESTAMP
        WHERE
			id = ?`,
		id,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("update player %d last_checkin: %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected check for player %d: %w", id, err)
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	return r.FindById(ctx, id)
}

//goland:noinspection
