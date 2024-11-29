package v1

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
)

type Item struct {
	DatabaseID byte
	Key        byte
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
		logger.Info("get item request received")
		databaseID := r.PathValue("databaseId")
		key := r.PathValue("key")

		fmt.Println(databaseID)
		fmt.Println(key)
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
