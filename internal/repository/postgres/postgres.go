package postgres

import (
	"context"
	"news-service/init/logger"
	"news-service/pkg/constants"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"

	"news-service/init/config"
)

func NewPostgresConnection(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	logger.Info("successfully connect to database. Migrating...", constants.LoggerPostgres)

	if err := goose.Up(db.DB, "./schema"); err != nil {

		return nil, err
	}

	logger.Info("successfully applying migrations", constants.LoggerPostgres)

	return db, nil
}
