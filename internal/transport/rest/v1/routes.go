package v1

import (
	"database/sql"
	"log/slog"
	"net/http"
)

func AttachSubRouter(router *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	subrouter := http.NewServeMux()

	subrouter.HandleFunc(`GET /read`, NewReadHandler(db, logger)) // id= and value=

	router.Handle(`/v1/`, http.StripPrefix(`/v1`, subrouter))
}
