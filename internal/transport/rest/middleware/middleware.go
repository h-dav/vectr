package middleware

import (
	"fmt"
	"net/http"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Authorization"))
		next.ServeHTTP(w, r)
	})
}
