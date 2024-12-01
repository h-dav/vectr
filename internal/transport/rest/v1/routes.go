package v1

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/h-dav/vectr/internal/transport/rest/middleware"
)

func AttachSubRouter(router *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	subrouter := http.NewServeMux()

	subrouter.Handle(`GET /filter`, middleware.Authentication(NewFilterHandler(db, logger)))
	subrouter.Handle(`POST /create`, middleware.Authentication(NewCreateHandler(db, logger)))

	router.Handle(`/v1/`, http.StripPrefix(`/v1`, subrouter))
}
