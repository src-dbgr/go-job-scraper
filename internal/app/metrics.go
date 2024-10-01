package app

import (
	"job-scraper/internal/metrics"
	"job-scraper/internal/storage"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func initJobCollector(storage storage.Storage) {
	jobCollector := metrics.NewJobCollector(storage)
	prometheus.MustRegister(jobCollector)
}

func startPrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Msgf("Starting Prometheus metrics server on :%d", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Error().Err(err).Msg("Prometheus metrics server failed")
	}
}
