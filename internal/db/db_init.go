package db

import (
	"Go-lab/internal/utils"
	"context"
	"database/sql"
	"log"
)

//goland:noinspection SqlNoDataSourceInspection
func createTables(ctx context.Context, dbUtils *utils.DbUtils) error {
	log.Println("Creating the table...")

	err := dbUtils.WithTransaction(func(tx *sql.Tx) error {
		var sqlString string
		var err error

		sqlString = "DROP TABLE IF EXISTS [client_entity];"
		_, err = dbUtils.DB.ExecContext(ctx, sqlString)
		if err != nil {
			return err
		}

		sqlString = `CREATE TABLE IF NOT EXISTS [client_entity] (
			[id] INTEGER PRIMARY KEY AUTOINCREMENT,
			[account_no] TEXT NOT NULL UNIQUE,
			[account_name] TEXT NOT NULL,
			[created_at] DATETIME DEFAULT CURRENT_TIMESTAMP);`

		_, err = dbUtils.DB.ExecContext(ctx, sqlString)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	log.Println("Created the table.")

	return nil
}

func PopulateTables() {

}
