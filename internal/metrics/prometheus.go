package metrics

import (
	// Re-export the metrics from domains
	"job-scraper/internal/metrics/domains"
)

var (
	ScrapedJobs     = domains.ScrapedJobs
	ProcessedJobs   = domains.ProcessedJobs
	ScraperDuration = domains.ScraperDuration
)

// Re-export the helper functions
func RecordScrapedJob()               { domains.RecordScrapedJob() }
func RecordProcessedJob()             { domains.RecordProcessedJob() }
func RecordScraperDuration(d float64) { domains.RecordScraperDuration(d) }
