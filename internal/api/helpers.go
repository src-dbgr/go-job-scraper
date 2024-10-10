package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
