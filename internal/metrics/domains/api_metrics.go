package domains

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "jobscraper",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "Duration of HTTP requests",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"handler", "method", "status"},
	)

	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "jobscraper",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"handler", "method", "status"},
	)

	ActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "jobscraper",
			Subsystem: "http",
			Name:      "active_requests",
			Help:      "Number of currently active requests",
		},
	)
)
