package middleware

import (
	"net/http"
	"strconv"
	"time"

	"job-scraper/internal/metrics/domains"

	"github.com/gorilla/mux"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		domains.ActiveRequests.Inc()
		defer domains.ActiveRequests.Dec()

		// Capture the response
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		route := mux.CurrentRoute(r)
		path := "unknown"
		if route != nil {
			pathTemplate, _ := route.GetPathTemplate()
			if pathTemplate != "" {
				path = pathTemplate
			}
		}

		statusCode := strconv.Itoa(rw.statusCode)

		domains.HTTPRequestDuration.WithLabelValues(
			path,
			r.Method,
			statusCode,
		).Observe(duration)

		domains.HTTPRequestsTotal.WithLabelValues(
			path,
			r.Method,
			statusCode,
		).Inc()
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
