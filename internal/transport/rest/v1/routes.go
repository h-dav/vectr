package v1

import (
	"database/sql"
	"log/slog"
	"net/http"
)

func AttachSubRouter(router *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	subrouter := http.NewServeMux()

	subrouter.HandleFunc(`POST /create-db`, NewCreateDbHandler(db, logger))
	subrouter.HandleFunc(`GET /read/{vectorId}`, NewReadHandler(db, logger))
	// subrouter.HandleFunc(`GET /read/{databaseId}/{value}`, NewReadHandler(db, logger))
	//subrouter.HandleFunc("GET /read", NewReadHandler(db, logger)) // A query - req.Body
	subrouter.HandleFunc(`POST /write`, NewWriteHandler(db, logger))
	subrouter.HandleFunc(`POST /update`, NewUpdateHandler(db, logger))
	subrouter.HandleFunc(`POST /delete`, NewDeleteHandler(db, logger))

	router.Handle(`/v1/`, http.StripPrefix(`/v1`, subrouter))
}
