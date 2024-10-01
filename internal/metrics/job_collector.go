package metrics

import (
	"context"
	"job-scraper/internal/storage"

	"github.com/prometheus/client_golang/prometheus"
)

type JobCollector struct {
	storage        storage.Storage
	jobsByCategory *prometheus.GaugeVec
	totalJobs      prometheus.Gauge
}

func NewJobCollector(storage storage.Storage) *JobCollector {
	return &JobCollector{
		storage: storage,
		jobsByCategory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "jobs_by_category",
				Help: "Number of jobs by category",
			},
			[]string{"category"},
		),
		totalJobs: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_jobs",
			Help: "Total number of jobs",
		}),
	}
}

func (c *JobCollector) Describe(ch chan<- *prometheus.Desc) {
	c.jobsByCategory.Describe(ch)
	c.totalJobs.Describe(ch)
}

func (c *JobCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	categoryCount, err := c.storage.GetJobCountByCategory(ctx)
	if err == nil {
		for category, count := range categoryCount {
			c.jobsByCategory.WithLabelValues(category).Set(float64(count))
		}
	}

	totalCount, err := c.storage.GetTotalJobCount(ctx)
	if err == nil {
		c.totalJobs.Set(float64(totalCount))
	}

	c.jobsByCategory.Collect(ch)
	c.totalJobs.Collect(ch)
}
