package domains

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Existing metrics from prometheus.go
var (
	ScrapedJobs = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "jobscraper",
		Subsystem: "basic",
		Name:      "scraped_jobs_total",
		Help:      "The total number of scraped jobs",
	})

	ProcessedJobs = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "jobscraper",
		Subsystem: "basic",
		Name:      "processed_jobs_total",
		Help:      "The total number of processed jobs",
	})

	ScraperDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "jobscraper",
		Subsystem: "basic",
		Name:      "scraper_duration_seconds",
		Help:      "Duration of scraper execution",
		Buckets:   prometheus.DefBuckets,
	})
)

// Export the existing functions to maintain compatibility
func RecordScrapedJob() {
	ScrapedJobs.Inc()
}

func RecordProcessedJob() {
	ProcessedJobs.Inc()
}

func RecordScraperDuration(duration float64) {
	ScraperDuration.Observe(duration)
}
