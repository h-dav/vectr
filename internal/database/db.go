package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/h-dav/vectr/internal/config"
	_ "github.com/lib/pq"
)

func NewConnection(cfg config.Config, logger *slog.Logger) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("opening DB connection", slog.Any("err", err)) // TODO: Find an alternative to slog.Any for errors.
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Error("pinging DB", slog.Any("err", err))
		return nil, err
	}

	return db, nil
}
