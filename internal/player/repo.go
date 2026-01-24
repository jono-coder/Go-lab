package player

import (
	"Go-lab/internal/utils/paging"
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	tx *sqlx.Tx
}

func NewRepo(tx *sqlx.Tx) (*Repo, error) {
	if err := validate.Get().Var(tx, "required"); err != nil {
		return nil, fmt.Errorf("invalid tx: %w", err)
	}

	return &Repo{tx: tx}, nil
}

func (r *Repo) Create(ctx context.Context, player *Player) (*uint, error) {
	if err := validate.Get().Var(player, "required"); err != nil {
		return nil, err
	}

	if err := player.Validate(); err != nil {
		return nil, err
	}

	res, err := r.tx.NamedExecContext(ctx,
		`INSERT INTO player_entity (resource_id, name, description) VALUES (:resource_id, :name, :description)`,
		&player,
	)
	if err != nil {
		return nil, fmt.Errorf("insert player: %w", err)
	}

	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("insert player (cannot get lastInsertId): %w", err)
	}

	id := uint(lastInsertedId) // so annoying, cannot be inlined

	return &id, nil
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func (r *Repo) FindById(ctx context.Context, id uint) (*Player, error) {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return nil, err
	}

	var player Player

	if err := r.tx.GetContext(ctx, &player, `
		SELECT
        	id,
			resource_id,
			name,
			description,
			last_checkin,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM
			player_entity
        WHERE
			id = ?
		AND
            deleted_at IS NULL`,
		id); err != nil {
		return nil, err
	}

	return &player, nil
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func (r *Repo) FindByResourceId(ctx context.Context, resourceId string) (*Player, error) {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return nil, err
	}
	if err := validate.Get().Var(resourceId, "notblank"); err != nil {
		return nil, err
	}

	var player Player

	if err := r.tx.GetContext(ctx, &player,
		`
		SELECT
			id,
			resource_id,
			name,
			description,
			last_checkin,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM
			player_entity
        WHERE
			resource_id = ?
		AND
            deleted_at IS NULL`,
		resourceId,
	); err != nil {
		return nil, err
	}

	return &player, nil
}

func (r *Repo) FindAll(ctx context.Context, paging paging.Paging) ([]Player, error) {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return nil, err
	}
	var players []Player

	if err := r.tx.SelectContext(ctx, &players,
		`
		SELECT
			id,
			resource_id,
			name,
			description,
			last_checkin,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM
			player_entity
		WHERE
		    deleted_at IS NULL
		ORDER BY
			name ASC
			LIMIT ? OFFSET ?`,
		paging.Limit, paging.Offset(),
	); err != nil {
		return nil, err
	}

	return players, nil
}

func (r *Repo) Checkin(ctx context.Context, id uint, updatedAt *time.Time) (*Player, error) {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return nil, err
	}

	res, err := r.tx.ExecContext(ctx, `
		UPDATE
			player_entity
		SET
			last_checkin = CURRENT_TIMESTAMP
		WHERE
			id = ?
		AND
			updated_at <=> ?`,
		id, updatedAt,
	)
	if err != nil {
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

func (r *Repo) Update(ctx context.Context, dto *UpdateDto) error {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return err
	}

	if err := validate.Get().Var(dto, "required"); err != nil {
		return err
	}

	if dto.Id == nil {
		return fmt.Errorf("id is required")
	}

	res, err := r.tx.NamedExecContext(ctx, `
		UPDATE
			player_entity
		SET
			name = :name,
			description = :description
		WHERE
			id = :id
		AND
			updated_at <=> :updated_at`,
		&dto,
	)
	if err != nil {
		return fmt.Errorf("update player %d: %w", *dto.Id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected check for player %d: %w", *dto.Id, err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete Soft Deletes only!
func (r *Repo) Delete(ctx context.Context, id uint, updatedAt *time.Time) error {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return err
	}

	res, err := r.tx.ExecContext(ctx, `
		UPDATE
			player_entity
		SET
			deleted_at = CURRENT_TIMESTAMP
		WHERE
			id = ?
		AND
			updated_at <=> ?
		AND
			deleted_at IS NULL`,
		id, updatedAt,
	)
	if err != nil {
		return fmt.Errorf("delete player %d: %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected check for player %d: %w", id, err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
