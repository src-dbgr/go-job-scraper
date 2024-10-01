package scraper

import (
	"context"
	"job-scraper/internal/models"
)

type Scraper interface {
	Scrape(ctx context.Context) ([]models.Job, error)
	Name() string
}
