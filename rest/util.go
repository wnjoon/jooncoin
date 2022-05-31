package rest

import (
	"fmt"
	"net/http"
)

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func setPort(_port int) string {
	return fmt.Sprintf(":%d", _port)
}

// Create middleware for add http header of json
// Called adapter using http.HandleFunc
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}
