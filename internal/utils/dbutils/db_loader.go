package dbutils

import (
	"Go-lab/internal/utils/session"
	"Go-lab/internal/utils/validate"
	"context"
	"database/sql"
	"log/slog"
	"os"
)

type DbLoader struct {
	db    *sql.DB
	utils *DbUtils
	ctx   context.Context
}

func NewDbLoader(ctx context.Context, db *sql.DB, dbUtils *DbUtils) *DbLoader {
	if err := validate.Required("ctx", ctx); err != nil {
		panic(err)
	}
	if err := validate.Required("db", db); err != nil {
		panic(err)
	}
	if err := validate.Required("utils", dbUtils); err != nil {
		panic(err)
	}

	return &DbLoader{
		db:    db,
		utils: dbUtils,
		ctx:   ctx,
	}
}

func (db *DbLoader) Load(ctx context.Context, scriptFilename string) error {
	scripts, err := os.ReadFile(scriptFilename)
	if err != nil {
		return err
	}

	slog.Info("running scripts...")

	ctx = session.ContextWithUserID(ctx, -1)

	err = db.utils.WithTransaction(ctx, func(context.Context, *sql.Tx) error {
		_, err := db.utils.DB.ExecContext(db.ctx, string(scripts))
		return err
	})

	slog.Info("ran scripts.")

	return err
}
