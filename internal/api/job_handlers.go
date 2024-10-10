package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (a *API) getJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50 // Default limit
	}

	jobs, err := a.storage.GetJobs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get jobs")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(jobs) > limit {
		jobs = jobs[:limit]
	}

	respondJSON(w, jobs)
}

func (a *API) getJobByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	job, err := a.storage.GetJobByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to get job")
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	respondJSON(w, job)
}

func (a *API) getJobUrls(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	urls, err := a.storage.GetExistingURLs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job URLs")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondJSON(w, urls)
}

func (a *API) getJobCategoryStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := a.storage.GetJobCountByCategory(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job category stats")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (a *API) getTotalJobCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	count, err := a.storage.GetTotalJobCount(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get total job count")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"total_jobs": count})
}
