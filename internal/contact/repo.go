package contact

import (
	"Go-lab/internal/utils/dbutils"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repo struct {
	db *dbutils.DbUtils
}

func NewRepo(dbUtils *dbutils.DbUtils) *Repo {
	return &Repo{
		db: dbUtils,
	}
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func (r *Repo) FindById(ctx context.Context, id int) (*Contact, error) {
	var res Contact

	err := r.db.DB.QueryRowContext(ctx,
		`SELECT id, first_name, surname, created_at
             	   FROM contact_entity
             	   WHERE id = ?`, id,
	).Scan(&res.Id, &res.firstName, &res.Surname, &res.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("find client %d: %w", id, err)
	}

	return &res, nil
}
