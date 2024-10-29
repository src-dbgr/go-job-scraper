package api

import (
	"context"
	"encoding/json"
	"job-scraper/internal/models"
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

	ctx := context.Background()

	existingURLs, err := a.storage.GetExistingURLs(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Str("scraper", scraperName).
			Msg("Failed to fetch existing URLs")
		a.runningScrapers.Store(scraperName, &ScraperStatus{Name: scraperName, Status: "Failed", Jobs: 0})
		return
	}

	var jobs []models.Job

	if paginatedScraper, ok := s.(scraper.PaginatedScraper); ok && pages > 0 {
		jobs, err = paginatedScraper.ScrapePages(ctx, pages)
	} else {
		jobs, err = s.Scrape(ctx)
	}

	if err != nil {
		log.Error().Err(err).Str("scraper", scraperName).Int("pages", pages).Msg("Failed to scrape jobs")
		a.runningScrapers.Store(scraperName, &ScraperStatus{Name: scraperName, Status: "Failed", Jobs: 0})
		return
	}

	processedJobs := 0
	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			if _, exists := existingURLs[job.URL]; !exists {
				processedJob, err := a.openaiProcessor.Process(ctx, job)
				if err != nil {
					log.Error().
						Err(err).
						Str("scraper", scraperName).
						Str("job_url", job.URL).
						Msg("Failed to process job with OpenAI")
					continue
				}

				if err := a.storage.SaveJob(ctx, processedJob); err != nil {
					log.Error().
						Err(err).
						Str("scraper", scraperName).
						Str("job_url", processedJob.URL).
						Msg("Failed to save processed job")
				} else {
					processedJobs++
					log.Info().
						Str("scraper", scraperName).
						Str("job_url", processedJob.URL).
						Str("job_title", processedJob.Title).
						Msg("Successfully processed and saved job")
				}
			} else {
				log.Info().Str("scraper", scraperName).Str("job_url", job.URL).Msg("Job already exists, skipping processing")
			}
		}
	}

	log.Info().
		Str("scraper", scraperName).
		Int("pages", pages).
		Int("total_jobs", len(jobs)).
		Int("processed_jobs", processedJobs).
		Msg("Scraping and processing completed")
	a.runningScrapers.Store(scraperName, &ScraperStatus{Name: scraperName, Status: "Completed", Jobs: processedJobs})
}
