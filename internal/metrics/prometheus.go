package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ScrapedJobs = promauto.NewCounter(prometheus.CounterOpts{
		Name: "scraped_jobs_total",
		Help: "The total number of scraped jobs",
	})

	ProcessedJobs = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processed_jobs_total",
		Help: "The total number of processed jobs",
	})

	ScraperDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "scraper_duration_seconds",
		Help:    "Duration of scraper execution",
		Buckets: prometheus.DefBuckets,
	})
)

func RecordScrapedJob() {
	ScrapedJobs.Inc()
}

func RecordProcessedJob() {
	ProcessedJobs.Inc()
}

func RecordScraperDuration(duration float64) {
	ScraperDuration.Observe(duration)
}
