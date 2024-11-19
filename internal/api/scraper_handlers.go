package api

import (
	"context"
	"encoding/json"
	"job-scraper/internal/scraper"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (a *API) handleScrape(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scraperName := vars["scraper"]

	s, ok := a.scrapers[scraperName]
	if !ok {
		http.Error(w, "Scraper not found", http.StatusNotFound)
		return
	}

	pages := 1 // Default value
	if pagesStr := r.URL.Query().Get("pages"); pagesStr != "" {
		if p, err := strconv.Atoi(pagesStr); err == nil && p > 0 {
			pages = p
		}
	}

	go a.runScraper(scraperName, s, pages)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Scraping job started", "scraper": scraperName, "pages": strconv.Itoa(pages)})
}

func (a *API) handleScrapersStatus(w http.ResponseWriter, r *http.Request) {
	statuses := []ScraperStatus{}
	a.runningScrapers.Range(func(key, value interface{}) bool {
		if status, ok := value.(*ScraperStatus); ok {
			statuses = append(statuses, *status)
		}
		return true
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func (a *API) runScraper(scraperName string, s scraper.Scraper, pages int) {
	a.runningScrapers.Store(scraperName, &ScraperStatus{Name: scraperName, Status: "Running", Jobs: 0})
	defer a.runningScrapers.Delete(scraperName)

	log.Info().
		Str("scraper", scraperName).
		Int("pages", pages).
		Msg("Starting scraper")

	result, err := a.scraperService.ExecuteScraping(context.Background(), s, pages)
	if err != nil {
		log.Error().Err(err).Str("scraper", scraperName).Msg("Scraping failed")
		a.runningScrapers.Store(scraperName, &ScraperStatus{
			Name:   scraperName,
			Status: "Failed",
			Jobs:   0,
		})
		return
	}
	log.Info().Msg("Job processing completed successfully")
	a.runningScrapers.Store(scraperName, &ScraperStatus{
		Name:   scraperName,
		Status: result.Status,
		Jobs:   result.ProcessedJobs,
	})
}
