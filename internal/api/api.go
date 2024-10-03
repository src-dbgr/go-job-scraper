package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"job-scraper/internal/models"
	"job-scraper/internal/scraper"
	"job-scraper/internal/storage"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type API struct {
	router          *mux.Router
	scrapers        map[string]scraper.Scraper
	storage         storage.Storage
	runningScrapers *sync.Map
}

type ScraperStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Jobs   int    `json:"jobs"`
}

func NewAPI(scrapers map[string]scraper.Scraper, storage storage.Storage) *API {
	api := &API{
		router:          mux.NewRouter(),
		scrapers:        scrapers,
		storage:         storage,
		runningScrapers: &sync.Map{},
	}
	api.setupRoutes()
	return api
}

func (a *API) setupRoutes() {
	a.router.HandleFunc("/api/scrape/{scraper}", a.handleScrape).Methods("POST")
	a.router.HandleFunc("/api/scrapers/status", a.handleScrapersStatus).Methods("GET")
}

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

	go func() {
		a.runningScrapers.Store(scraperName, &ScraperStatus{Name: scraperName, Status: "Running", Jobs: 0})
		defer a.runningScrapers.Delete(scraperName)

		ctx := context.Background()
		var jobs []models.Job
		var err error

		if paginatedScraper, ok := s.(scraper.PaginatedScraper); ok && pages > 1 {
			jobs, err = paginatedScraper.ScrapePages(ctx, pages)
		} else {
			jobs, err = s.Scrape(ctx)
		}

		if err != nil {
			log.Error().Err(err).Str("scraper", scraperName).Int("pages", pages).Msg("Failed to scrape jobs")
			a.runningScrapers.Store(scraperName, &ScraperStatus{Name: scraperName, Status: "Failed", Jobs: 0})
			return
		}

		for _, job := range jobs {
			if err := a.storage.SaveJob(ctx, job); err != nil {
				log.Error().Err(err).Str("scraper", scraperName).Msg("Failed to save job")
			}
		}

		log.Info().Str("scraper", scraperName).Int("pages", pages).Int("jobs", len(jobs)).Msg("Scraping completed")
		a.runningScrapers.Store(scraperName, &ScraperStatus{Name: scraperName, Status: "Completed", Jobs: len(jobs)})
	}()

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

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
