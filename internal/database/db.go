package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/h-dav/vectr/internal/config"
	_ "github.com/lib/pq"
)

func NewConnection(cfg config.Config, logger *slog.Logger) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?ssl-mode=verify-full", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("opening DB connection: %w", slog.Any("err", err)) // TODO: Find an alternative to slog.Any for errors.
		return nil, err
	}

	return db, nil

}
