package session_db

import (
	"Go-lab/internal/utils/validate"
	"context"

	"github.com/jmoiron/sqlx"
)

func GetUserIdFromDB(ctx context.Context, tx *sqlx.Tx) (*int, error) {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return nil, err
	}
	if err := validate.Get().Var(tx, "required"); err != nil {
		return nil, err
	}

	var userID *int

	err := tx.QueryRowContext(ctx, "SELECT get_current_user_id()").Scan(&userID)
	if err != nil {
		return nil, err
	}

	if userID == nil {
		return nil, nil
	}

	return userID, err
}
