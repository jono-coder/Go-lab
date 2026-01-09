package session_db

import (
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
)

func GetUserIdFromDB(ctx context.Context, tx *sql.Tx) (int, error) {
	err := validate.Required("tx", tx)
	if err != nil {
		return 0, err
	}

	var userID *int

	err = tx.QueryRowContext(ctx, "SELECT get_current_user_id()").Scan(&userID)
	if err != nil {
		return -1, err
	}

	if userID == nil {
		return -1, nil
	}

	return *userID, err
}
