package v1

import (
	"database/sql"
	"log/slog"
	"net/http"
)

type Item struct {
	DatabaseID byte
	ID         byte
	Value      byte
	Vector     []float64
	Metadata   map[any]any
}

func NewCreateDbHandler(db *sql.DB, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("post request received")
	}
}
func NewReadHandler(db *sql.DB, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		databaseID := r.PathValue("databaseId")

		// key := r.PathValue("key")

		rows, err := db.Query(`select value from vectors where database_id = $1`, databaseID)
		if err != nil {
			logger.Error("connection query error", slog.Any("err", err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		for rows.Next() {
			var value string
			if err := rows.Scan(&value); err != nil {
				logger.Error("error: %w", slog.Any("err", err))
				w.WriteHeader(http.StatusInternalServerError)
			}

			w.Write([]byte(value))
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}

func NewWriteHandler(db *sql.DB, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("create request received")
	}
}

func NewUpdateHandler(db *sql.DB, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("update request received")
	}
}

func NewDeleteHandler(db *sql.DB, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("delete request received")
	}
}
