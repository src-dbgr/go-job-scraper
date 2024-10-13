package api

import (
	"net/http"
)

// VersionMiddleware adds the API version to the response header
func VersionMiddleware(next http.HandlerFunc, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("API-Version", version)
		next.ServeHTTP(w, r)
	}
}
