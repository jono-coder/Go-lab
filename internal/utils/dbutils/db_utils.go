package dbutils

import (
	"Go-lab/config"
	"Go-lab/internal/utils/session"
	"context"
	"database/sql"
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

func (dbUtils *DbUtils) WithTransaction(ctx context.Context, txFunc func(*sql.Tx) error) error {
	con, err := dbUtils.DB.Conn(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := resetSession(ctx, con); err != nil {
			slog.Error("failed to reset session", "error", err)
		}
		if err := con.Close(); err != nil {
			slog.Error("failed to close connection", "error", err)
		}
	}()

	tx, err := con.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = initSessionVars(ctx, tx)
	if err != nil {
		return err
	}

	if err := txFunc(tx); err != nil {
		return err
	}

	return tx.Commit()
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

func resetSession(ctx context.Context, con *sql.Conn) error {
	_, err := con.ExecContext(ctx, "SET @session_user_id = NULL")
	return err
}

func initSessionVars(ctx context.Context, tx *sql.Tx) error {
	userId, found := session.UserIDFromContext(ctx)
	if !found {
		slog.Warn("No user ID found in context!")
		userId = -1
	}

	_, err := tx.ExecContext(ctx, "SET @session_user_id = ?", userId)
	return err
}
