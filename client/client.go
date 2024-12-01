package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Vector struct {
	ID         byte   `json:"vector_id"`
	DatabaseID byte   `json:"database_id"`
	Value      string `json:"value"`
	Vector     []byte `json:"vector"`
	Metadata   string `json:"metadata"`
}

func GetVectorByID(id string) (Vector, error) {
	requestURL := fmt.Sprintf("http://localhost:8080/v1/read/1")
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return Vector{}, nil
	}

	var vector Vector
	json.NewDecoder(res.Body).Decode(&vector)

	return vector, nil

}
