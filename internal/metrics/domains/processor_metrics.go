package domains

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProcessorDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "jobscraper",
			Subsystem: "processor",
			Name:      "operation_duration_seconds",
			Help:      "Duration of processing operations",
			Buckets:   []float64{0.5, 1, 2, 5, 10, 20},
		},
		[]string{"processor", "status"},
	)

	ProcessorErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "jobscraper",
			Subsystem: "processor",
			Name:      "errors_total",
			Help:      "Total number of processor errors",
		},
		[]string{"processor", "error_type"},
	)

	OpenAITokensUsed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "jobscraper",
			Subsystem: "processor",
			Name:      "openai_tokens_total",
			Help:      "Total number of OpenAI tokens used",
		},
		[]string{"model", "operation"},
	)
)
