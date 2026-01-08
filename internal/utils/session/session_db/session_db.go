package session_db

import (
	"Go-lab/internal/utils/dbutils"
	"context"
	"database/sql"
)

func GetUserIdFromDB(ctx context.Context, db *dbutils.DbUtils) (int, error) {
	var userID int

	err := db.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, "SELECT get_current_user_id()").Scan(&userID)
	})

	if err != nil {
		return -1, err
	}

	return userID, err
}
