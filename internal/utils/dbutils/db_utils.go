package utils

import (
	"Go-lab/config"
	"Go-lab/internal/utils/session"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	//_ "modernc.org/sqlite"
)

type DbUtils struct {
	DB     *sql.DB
	config *config.DBConfig
	mtx    sync.Mutex
}

func NewDbUtils(config *config.DBConfig) *DbUtils {
	res := &DbUtils{
		DB:     open(config),
		config: config,
	}
	return res
}

func (dbUtils *DbUtils) Init() error {
	dbUtils.mtx.Lock()
	defer dbUtils.mtx.Unlock()

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = dbUtils.createTables(ctx)
	if err != nil {
		return err
	}

	err = dbUtils.populateTables(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (dbUtils *DbUtils) WithTransaction(ctx context.Context, txFunc func(context.Context, *sql.Tx) error) error {
	tx, err := dbUtils.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if userId, found := session.UserIDFromContext(ctx); found {
		_, err := tx.ExecContext(ctx, "SET @session_user_id = ?", userId)
		if err != nil {
			return err
		}
	} else {
		slog.Warn("No user ID found in context!")
	}

	if err := txFunc(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (dbUtils *DbUtils) Close() {
	log.Println("Closing the database...")

	dbUtils.mtx.Lock()
	defer dbUtils.mtx.Unlock()

	if dbUtils.DB == nil {
		return
	}

	err := dbUtils.DB.Close()
	if err != nil {
		log.Printf("Couldn't close the database :: %v", err)
	}

	dbUtils.DB = nil

	log.Println("Closed the database.")
}

func open(config *config.DBConfig) *sql.DB {
	log.Println("Opening the database...")

	// Open the database
	res, err := sql.Open(config.Driver, config.DSN)
	if err != nil {
		panic(err)
	}

	res.SetConnMaxLifetime(time.Minute * 3)
	res.SetMaxOpenConns(10)
	res.SetMaxIdleConns(10)

	log.Println("Opened the database.")

	return res
}

//goland:noinspection SqlNoDataSourceInspection
func (dbUtils *DbUtils) createTables(ctx context.Context) error {
	err := dbUtils.WithTransaction(func(tx *sql.Tx) error {
		var sqlString string
		var err error

		log.Println("Creating client_entity table...")

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

		log.Println("Created client_entity table.")

		log.Println("Creating player_entity table...")

		sqlString = "DROP TABLE IF EXISTS [player_entity];"
		_, err = dbUtils.DB.ExecContext(ctx, sqlString)
		if err != nil {
			return err
		}

		sqlString = `CREATE TABLE IF NOT EXISTS [player_entity] (
			[id] INTEGER PRIMARY KEY AUTOINCREMENT,
			[resource_id] TEXT,
			[name] TEXT NOT NULL,
			[description] TEXT,
			[last_checkin] DATETIME,
			[created_at] DATETIME DEFAULT CURRENT_TIMESTAMP);`

		_, err = dbUtils.DB.ExecContext(ctx, sqlString)
		if err != nil {
			return err
		}

		log.Println("Created player_entity table.")


		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection


type Data struct {
    Clients [][]string
    Players [][]string
}

func (dbUtils *DbUtils) populateTables(ctx context.Context) error {

	data := Data{
		Clients: [][]string{
			{"ABC001", "ABC Shoes"},
			{"XYZ001", "XYZ Trading As XYZ Enterprises"},
		},
		Players: [][]string{
			{"abcd1234", "Player One", "1st example player"},
			{"defg5678", "Player Two", "2nd example player"},
		},
func (dbUtils *DbUtils) populateTablesX(ctx context.Context) error {
	data := [][]string{
		{"ABC001", "ABC Shoes"},
		{"XYZ001", "XYZ Trading As XYZ Enterprises"},
	}

	clientStmt, err := dbUtils.DB.PrepareContext(ctx, `INSERT INTO client_entity (account_no, account_name) VALUES (?, ?);`)
	if err != nil { return err }
	defer clientStmt.Close()

	playerStmt, err := dbUtils.DB.PrepareContext(ctx, `INSERT INTO player_entity (resource_id, name, description) VALUES (?, ?, ?);`)
	if err != nil { return err }
	defer playerStmt.Close()

	log.Println("Populating client_entity table...")

	for _, row := range data.Clients {
		if len(row) < 2 { return fmt.Errorf("invalid client row: %#v", row) }
		res, err := clientStmt.ExecContext(ctx, row[0], row[1])
		if err != nil { return fmt.Errorf("insert client failed (code=%s): %w", row[0], err) }
		if n, _ := res.RowsAffected(); n == 0 { log.Printf("client not inserted (code=%s)", row[0]) }
	}

	log.Println("Populated client_entity table.")

	log.Println("Populating player_entity table...")

	for _, row := range data.Players {
		if len(row) < 3 { return fmt.Errorf("invalid player row: %#v", row) }
		res, err := playerStmt.ExecContext(ctx, row[0], row[1], row[2])
		if err != nil { return fmt.Errorf("insert player failed (code=%s): %w", row[0], err) }
		if n, _ := res.RowsAffected(); n == 0 { log.Printf("player not inserted (code=%s)", row[0]) }
	}

	log.Println("Populated player_entity table.")

	return nil
}
