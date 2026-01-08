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
