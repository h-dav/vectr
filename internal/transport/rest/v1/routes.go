package v1

import (
	"database/sql"
	"log/slog"
	"net/http"
)

func AttachSubRouter(router *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	subrouter := http.NewServeMux()

	subrouter.HandleFunc(`GET /filter`, NewFilterHandler(db, logger))
	subrouter.HandleFunc(`POST /create`, NewCreateHandler(db, logger))

	router.Handle(`/v1/`, http.StripPrefix(`/v1`, subrouter))
}
