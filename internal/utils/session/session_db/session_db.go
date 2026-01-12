package session_db

import (
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
)

func GetUserIdFromDB(ctx context.Context, tx *sql.Tx) (*int, error) {
	if err := validate.Required("ctx", ctx); err != nil {
		return nil, err
	}
	if err := validate.Required("tx", tx); err != nil {
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
