package domains

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ScrapingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "jobscraper",
			Subsystem: "scraper",
			Name:      "operation_duration_seconds",
			Help:      "Duration of scraping operations",
			Buckets:   []float64{0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"scraper", "status"},
	)

	ScrapedJobsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "jobscraper",
			Subsystem: "scraper",
			Name:      "jobs_total",
			Help:      "Total number of scraped jobs",
		},
		[]string{"scraper", "status"},
	)

	ScraperErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "jobscraper",
			Subsystem: "scraper",
			Name:      "errors_total",
			Help:      "Total number of scraper errors",
		},
		[]string{"scraper", "error_type"},
	)
)
