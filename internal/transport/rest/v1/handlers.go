package v1

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

type Vec []float64
type Metadata map[string]any

type Vector struct {
	ID         byte     `json:"vector_id"`
	DatabaseID byte     `json:"database_id"`
	Value      string   `json:"value"`
	Vector     Vec      `json:"vector"`
	Metadata   Metadata `json:"metadata"`
}

func NewFilterHandler(db *sql.DB, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("filter request received")

		vectorId := r.URL.Query().Get("id")

		if vectorId != "" {
			vector, err := getVectorByID(db, vectorId, logger)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			if vector.ID == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			data, _ := json.Marshal(vector)

			json.Unmarshal(data, &vector)

			w.WriteHeader(http.StatusOK)
			w.Write(data)

			return
		}

		value := r.URL.Query().Get("value")

		if value != "" {
			vectors, err := getVectorsByValue(db, value, logger)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			if len(vectors) == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			data, _ := json.Marshal(vectors)

			json.Unmarshal(data, &vectors)

			w.WriteHeader(http.StatusOK)
			w.Write(data)

			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func getVectorsByValue(db *sql.DB, value string, logger *slog.Logger) ([]Vector, error) {
	rows, err := db.Query(`select vector_id, database_id, value, vector, metadata from vectors where value = $1`, value)
	if err != nil {
		logger.Error("connection query error", slog.Any("err", err))
		return []Vector{}, err
	}
	defer rows.Close()

	var vectors []Vector

	for rows.Next() {
		var (
			vector Vector
			v      []byte
			md     []byte
		)

		if err := rows.Scan(&vector.ID, &vector.DatabaseID, &vector.Value, &v, &md); err != nil {
			logger.Error("scanning row", slog.Any("err", err))
			return []Vector{}, err
		}

		vector.Vector, _ = convertBytesToVec(v, logger)
		vector.Metadata, _ = convertBytesToMetadata(md, logger)

		vectors = append(vectors, vector)
	}

	return vectors, nil
}

func getVectorByID(db *sql.DB, vectorID string, logger *slog.Logger) (Vector, error) {
	rows, err := db.Query(`select vector_id, database_id, value, vector, metadata from vectors where vector_id = $1`, vectorID)
	if err != nil {
		logger.Error("connection query error", slog.Any("err", err))
		return Vector{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			vector Vector
			v      []byte
			md     []byte
		)

		if err := rows.Scan(&vector.ID, &vector.DatabaseID, &vector.Value, &v, &md); err != nil {
			logger.Error("scanning row", slog.Any("err", err))
			return Vector{}, err
		}

		vector.Vector, _ = convertBytesToVec(v, logger)
		vector.Metadata, _ = convertBytesToMetadata(md, logger)

		return vector, nil // This could lead to bugs further down the line in the case that two vectors are returned.
	}

	return Vector{}, nil
}

func NewCreateHandler(db *sql.DB, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("create request received")

		var vector Vector

		if err := json.NewDecoder(r.Body).Decode(&vector); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err := insertVector(db, vector, logger); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func insertVector(db *sql.DB, vector Vector, logger *slog.Logger) (string, error) {
	md, _ := convertMetadataToBytes(vector.Metadata, logger)

	sql := "insert into vectors (vector_id, database_id, value, vector, metadata) values ($1, $2, $3, $4, $5)"
	if _, err := db.Exec(sql, vector.ID, vector.DatabaseID, vector.Value, pq.Float64Array(vector.Vector), md); err != nil {
		logger.Error("connection query error", slog.Any("err", err))
		return "", err
	}

	return string(vector.ID), nil
}

func convertBytesToVec(bytes []byte, logger *slog.Logger) (Vec, error) {
	str := string(bytes)
	str = strings.Trim(str, "{")
	str = strings.Trim(str, "}")

	parts := strings.Split(str, ",")

	var vector Vec

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

func convertBytesToMetadata(bytes []byte, logger *slog.Logger) (Metadata, error) {
	var metadata Metadata

	if err := json.Unmarshal(bytes, &metadata); err != nil {
		logger.Error("parsing metadata", slog.Any("err", err))
		return Metadata{}, err
	}

	return metadata, nil
}

func convertMetadataToBytes(metadata Metadata, logger *slog.Logger) ([]byte, error) {
	bytes, err := json.Marshal(metadata)
	if err != nil {
		logger.Error("parsing metadata", slog.Any("err", err))
		return []byte{}, err
	}

	return bytes, nil
}
