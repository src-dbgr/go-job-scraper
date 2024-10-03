package scraper

import (
	"context"
	"job-scraper/internal/models"
)

// Scraper ist das Basis-Interface für alle Scraper
type Scraper interface {
	Scrape(ctx context.Context) ([]models.Job, error)
	Name() string
}

// PaginatedScraper ist ein erweitertes Interface für Scraper mit Paging-Unterstützung
type PaginatedScraper interface {
	Scraper
	ScrapePages(ctx context.Context, pages int) ([]models.Job, error)
}
