package services

import (
	"context"
	"job-scraper/internal/models"
	"job-scraper/internal/processor"
	"job-scraper/internal/scraper"
	"job-scraper/internal/storage"

	"github.com/rs/zerolog/log"
)

// ScraperService kapselt die Scraping-Logik
type ScraperService struct {
	storage   storage.Storage
	processor processor.JobProcessor
}

// ScrapingResult repräsentiert das Ergebnis eines Scraping-Durchlaufs
type ScrapingResult struct {
	TotalJobs     int
	ProcessedJobs int
	Status        string
	Error         error
}

// NewScraperService erstellt eine neue Instanz des ScraperService
func NewScraperService(storage storage.Storage, processor processor.JobProcessor) *ScraperService {
	return &ScraperService{
		storage:   storage,
		processor: processor,
	}
}

// ExecuteScraping führt den vollständigen Scraping-Workflow aus
func (s *ScraperService) ExecuteScraping(ctx context.Context, scraper scraper.Scraper, pages int) (*ScrapingResult, error) {
	result := &ScrapingResult{
		Status: "Running",
	}

	existingURLs, getExistingUrlError := s.storage.GetExistingURLs(ctx)
	if getExistingUrlError != nil {
		log.Error().Err(getExistingUrlError).Msg("Failed to fetch existing URLs")
		result.Status = "Failed"
		result.Error = getExistingUrlError
		return result, getExistingUrlError
	}

	var jobs []models.Job
	var err error

	// Type assertion für PaginatedScraper
	if pScraper, ok := scraper.(interface {
		ScrapePages(ctx context.Context, pages int) ([]models.Job, error)
	}); ok && pages > 0 {
		jobs, err = pScraper.ScrapePages(ctx, pages)
	} else {
		jobs, err = scraper.Scrape(ctx)
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to scrape jobs")
		result.Status = "Failed"
		result.Error = err
		return result, err
	}

	result.TotalJobs = len(jobs)

	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
			if err := s.processJob(ctx, job, existingURLs, result); err != nil {
				log.Error().Err(err).Str("job_url", job.URL).Msg("Failed to process job")
			}
		}
	}

	result.Status = "Completed"
	return result, nil
}

func (s *ScraperService) processJob(ctx context.Context, job models.Job, existingURLs map[string]bool, result *ScrapingResult) error {
	if _, exists := existingURLs[job.URL]; exists {
		log.Info().Str("job_url", job.URL).Msg("Job already exists, skipping processing")
		return nil
	}

	processedJob, err := s.processor.Process(ctx, job)
	if err != nil {
		return err
	}

	if err := s.storage.SaveJob(ctx, processedJob); err != nil {
		return err
	}

	result.ProcessedJobs++
	log.Info().
		Str("job_url", processedJob.URL).
		Str("job_title", processedJob.Title).
		Msg("Successfully processed and saved job")

	return nil
}
