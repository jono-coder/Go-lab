package dbutils

import (
	"Go-lab/config"
	"Go-lab/internal/utils/session"
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
	"log"
	"log/slog"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	//_ "modernc.org/sqlite"
)

type DbUtils struct {
	DB     *sqlx.DB
	config *config.DBConfig
	mtx    sync.Mutex
}

func NewDbUtils(config *config.DBConfig) *DbUtils {
	if err := validate.Get().Var(config, "required"); err != nil {
		panic(err)
	}

	return &DbUtils{
		DB:     open(config),
		config: config,
	}
}

func (dbUtils *DbUtils) WithTransaction(ctx context.Context, txFunc func(*sqlx.Tx) error) error {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return err

	}
	if err := validate.Get().Var(txFunc, "required"); err != nil {
		return err
	}

	con, err := dbUtils.DB.Connx(ctx)
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

	err = initSessionVars(ctx, con)
	if err != nil {
		return err
	}

	tx, err := con.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := txFunc(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func ToTime(time sql.NullTime) *time.Time {
	if time.Valid {
		return &time.Time
	}
	return nil
}

func open(config *config.DBConfig) *sqlx.DB {
	log.Println("Opening the database...")

	// Open the database
	db, err := sqlx.Open(config.Driver, config.DSN)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	log.Println("Opened the database.")

	return db
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

func resetSession(ctx context.Context, con *sqlx.Conn) error {
	_, err := con.ExecContext(ctx, "SET @session_user_id = NULL")
	return err
}

func initSessionVars(ctx context.Context, conn *sqlx.Conn) error {
	userId, found := session.UserIDFromContext(ctx)

	var err error
	if !found {
		slog.Warn("No user ID found in context!")
		_, err = conn.ExecContext(ctx, "SET @session_user_id = NULL")
	} else {
		_, err = conn.ExecContext(ctx, "SET @session_user_id = ?", userId)
	}

	return err
}
