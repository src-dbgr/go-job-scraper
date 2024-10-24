package domains

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	DBOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "jobscraper",
			Subsystem: "storage",
			Name:      "operation_duration_seconds",
			Help:      "Duration of database operations",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"operation", "status"},
	)

	DBOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "jobscraper",
			Subsystem: "storage",
			Name:      "operations_total",
			Help:      "Total number of database operations",
		},
		[]string{"operation", "status"},
	)

	DBConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "jobscraper",
			Subsystem: "storage",
			Name:      "connections_active",
			Help:      "Number of active database connections",
		},
	)
)
