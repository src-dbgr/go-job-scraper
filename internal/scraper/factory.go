package scraper

import (
	"fmt"
	"job-scraper/internal/models"
	"job-scraper/internal/scraper/jobsch"
	"net/http"
	"strconv"
)

func NewScraper(name string, config map[string]string) (Scraper, error) {
	switch name {
	case JobsChScraperName:
		return newJobschScraper(config)
	// to be extended
	default:
		return nil, fmt.Errorf("unknown scraper: %s", name)
	}
}

func newJobschScraper(config map[string]string) (Scraper, error) {
	baseURL, ok := config["base_url"]
	if !ok {
		return nil, fmt.Errorf("base_url is required for jobsch scraper")
	}

	maxPages, err := strconv.Atoi(config["max_pages"])
	if err != nil {
		maxPages = 10 // Default value if not specified or invalid
	}

	pageSize, err := strconv.Atoi(config["page_size"])
	if err != nil {
		pageSize = 20 // Default value if not specified or invalid
	}

	client := &http.Client{}
	fetcher := jobsch.NewJobsChFetcher(client, baseURL)

	scraperConfig := jobsch.Config{
		BaseURL:    baseURL,
		MaxPages:   maxPages,
		PageSize:   pageSize,
		JobFetcher: fetcher,
		ParseFunc: func(data []byte) (*models.Job, error) {
			// TODO: Implement proper parsing logic
			return &models.Job{}, nil
		},
	}

	return jobsch.NewJobsChScraper(scraperConfig), nil
}
