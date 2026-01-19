package dbutils

import (
	"Go-lab/internal/utils/session"
	"Go-lab/internal/utils/validate"
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
)

type DbLoader struct {
	utils *DbUtils        `validate:"required"`
	ctx   context.Context `validate:"required"`
}

func NewDbLoader(ctx context.Context, dbUtils *DbUtils) *DbLoader {
	res := &DbLoader{
		utils: dbUtils,
		ctx:   ctx,
	}
	err := validate.Get().Struct(res)
	if err != nil {
		panic(err)
	}
	return res
}

func (db *DbLoader) Load(ctx context.Context, scriptFilename string) error {
	if err := validate.Get().Var(ctx, "required"); err != nil {
		return err
	}
	if err := validate.NotBlank("scriptFilename", scriptFilename); err != nil {
		return err
	}

	scripts, err := os.ReadFile(scriptFilename)
	if err != nil {
		return err
	}

	slog.Info("running scripts...")

	// admin user
	ctx = session.ContextWithUserID(ctx, 0)

	err = db.utils.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		userID, found := session.UserIDFromContext(ctx)
		if found {
			log.Printf("user id: %d", userID)
		} else {
			log.Printf("user id is NULL")
		}
		_, err := tx.ExecContext(db.ctx, string(scripts))
		return err
	})

	slog.Info("ran scripts.")

	return err
}
