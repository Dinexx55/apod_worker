package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"nasa-apod-app/internal/config"
	"nasa-apod-app/internal/migration"
	"nasa-apod-app/internal/repository/postgres"
	"time"
)

func InitDB(cfg config.DBConfig, migrator *migration.Migratory, logger *zap.Logger) (repo *postgres.ApodImagesRepository, err error) {
	logger.Info("Got db config")

	var db *sqlx.DB
	for i := 0; i < cfg.ReconnRetry; i++ {

		db, err = postgres.ConnectToPostgresDB(cfg, logger)
		if err == nil {
			logger.Info("Db migration")

			if err = migrator.Migrate(db); err != nil {
				return nil, fmt.Errorf("migration failure: %w", err)
			}

			repo = postgres.NewPostgresRepository(db)
			logger.Info("Migrations done")

			return repo, nil
		}

		logger.With(
			zap.String("place", "main"),
			zap.Error(err),
		).Error("Failed to connect to db. Retrying")

		time.Sleep(cfg.TimeWaitPerTry)
	}

	logger.Info("Successfully connected to postgres")
	return repo, err
}
