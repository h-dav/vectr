package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type Vec []float64

type Vector struct {
	ID         byte   `json:"vector_id"`
	DatabaseID byte   `json:"database_id"`
	Value      string `json:"value"`
	Vector     Vec    `json:"vector"`
	Metadata   string `json:"metadata"`
}

type VectorResponse struct {
	ID         byte   `json:"vector_id"`
	DatabaseID byte   `json:"database_id"`
	Value      string `json:"value"`
	Vector     string `json:"vector"`
	Metadata   string `json:"metadata"`
}

func NewCreateDbHandler(db *sql.DB, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("post request received")
	}
}
func NewReadHandler(db *sql.DB, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vectorId := r.PathValue("vectorId")

		vector, err := getVectorsByIDs(db, vectorId, logger)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		logger.Info("value", slog.Any("vec", vector.Vector))

		data, _ := json.Marshal(vector)

		json.Unmarshal(data, &vector)

		fmt.Println(vector)

		fmt.Println(string(data))

		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
}

func getVectorsByIDs(db *sql.DB, vectorID string, logger *slog.Logger) (Vector, error) {
	rows, err := db.Query(`select vector_id, database_id, value, vector, metadata from vectors where vector_id = $1`, vectorID)
	if err != nil {
		logger.Error("connection query error", slog.Any("err", err))
		return Vector{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var vector Vector
		var v []byte

		if err := rows.Scan(&vector.ID, &vector.DatabaseID, &vector.Value, &v, &vector.Metadata); err != nil {
			logger.Error("scanning row", slog.Any("err", err))
			return Vector{}, err
		}

		vector.Vector, _ = convertBytesToVec(v, logger)

		return vector, nil // This could lead to bugs further down the line in the case that two vectors are returned.
	}

	return Vector{}, nil
}

func convertBytesToVec(bytes []byte, logger *slog.Logger) (Vec, error) {
	var vector Vec

	str := string(bytes)
	str = strings.Trim(str, "{")
	str = strings.Trim(str, "}")

	parts := strings.Split(str, ",")

	for _, part := range parts {
		f, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			logger.Error("parsing float", slog.Any("err", err))
			return Vec{}, nil
		}

		vector = append(vector, f)
	}

	return vector, nil
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
