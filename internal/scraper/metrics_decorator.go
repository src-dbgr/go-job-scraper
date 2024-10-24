package scraper

import (
	"context"
	"time"

	"job-scraper/internal/metrics/domains"
	"job-scraper/internal/models"
)

type MetricsDecorator struct {
	scraper Scraper
}

func NewMetricsDecorator(scraper Scraper) Scraper {
	return &MetricsDecorator{scraper: scraper}
}

func (d *MetricsDecorator) Scrape(ctx context.Context) ([]models.Job, error) {
	start := time.Now()
	jobs, err := d.scraper.Scrape(ctx)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
		domains.ScraperErrors.WithLabelValues(d.scraper.Name(), "scrape_error").Inc()
	}

	domains.ScrapingDuration.WithLabelValues(d.scraper.Name(), status).Observe(duration)
	domains.ScrapedJobsTotal.WithLabelValues(d.scraper.Name(), status).Add(float64(len(jobs)))

	return jobs, err
}

func (d *MetricsDecorator) Name() string {
	return d.scraper.Name()
}

// PaginatedMetricsDecorator implementiert auch das PaginatedScraper Interface
type PaginatedMetricsDecorator struct {
	MetricsDecorator
	scraper PaginatedScraper
}

func NewPaginatedMetricsDecorator(scraper PaginatedScraper) PaginatedScraper {
	return &PaginatedMetricsDecorator{
		MetricsDecorator: MetricsDecorator{scraper: scraper},
		scraper:          scraper,
	}
}

func (d *PaginatedMetricsDecorator) ScrapePages(ctx context.Context, pages int) ([]models.Job, error) {
	start := time.Now()
	jobs, err := d.scraper.ScrapePages(ctx, pages)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
		domains.ScraperErrors.WithLabelValues(d.scraper.Name(), "scrape_pages_error").Inc()
	}

	domains.ScrapingDuration.WithLabelValues(d.scraper.Name(), status).Observe(duration)
	domains.ScrapedJobsTotal.WithLabelValues(d.scraper.Name(), status).Add(float64(len(jobs)))

	return jobs, err
}
