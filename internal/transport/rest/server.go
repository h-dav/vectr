package rest

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/h-dav/vectr/internal/config"
	v1 "github.com/h-dav/vectr/internal/transport/rest/v1"
)

const readerHeaderTimeout = time.Second * 10

func NewServer(cfg config.Config, db *sql.DB, logger *slog.Logger) (*http.Server, error) {
	router := http.NewServeMux()

	attachHealthRouter(router, logger)

	v1.AttachSubRouter(router, db, logger)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTP.Port,
		ReadHeaderTimeout: readerHeaderTimeout,
		Handler:           router,
	}

	return srv, nil
}

func attachHealthRouter(router *http.ServeMux, logger *slog.Logger) {
	router.HandleFunc(`/alive`, func(w http.ResponseWriter, r *http.Request) {
		logger.Info("alive check received")
		w.WriteHeader(http.StatusOK)
	})
}
