package middleware

import (
	"encoding/json"
	"job-scraper/internal/apperrors"
	"net/http"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				handleError(w, apperrors.NewBaseError("INTERNAL_SERVER_ERROR", "Internal server error", nil))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	switch e := err.(type) {
	case *apperrors.NotFoundError:
		statusCode = http.StatusNotFound
		response.Code = e.Code
		response.Message = e.Message
	case *apperrors.ScrapingError:
		statusCode = http.StatusServiceUnavailable
		response.Code = e.Code
		response.Message = e.Message
	case *apperrors.ProcessingError:
		statusCode = http.StatusInternalServerError
		response.Code = e.Code
		response.Message = e.Message
	case *apperrors.BaseError:
		statusCode = http.StatusInternalServerError
		response.Code = e.Code
		response.Message = e.Message
	default:
		statusCode = http.StatusInternalServerError
		response.Code = apperrors.ErrCodeInternal
		response.Message = "An unexpected error occurred"
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
