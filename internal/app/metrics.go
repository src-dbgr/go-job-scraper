package app

import (
	"job-scraper/internal/metrics"
	"job-scraper/internal/metrics/domains"
	"job-scraper/internal/storage"
	"job-scraper/internal/storage/mongodb"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func initMetrics(storage storage.Storage) {
	// Register existing collectors
	jobCollector := metrics.NewJobCollector(storage)
	prometheus.MustRegister(jobCollector)

	// Initialize connection tracking
	go trackConnections(storage)
}

func trackConnections(storage storage.Storage) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Versuche, den konkreten MongoDB-Client zu erhalten
		if mongoClient, ok := storage.(*mongodb.Client); ok {
			if client := mongoClient.GetMongoClient(); client != nil {
				// Verwende den Mongo-Client für die Statistiken
				stats := float64(client.NumberSessionsInProgress())
				domains.DBConnectionsActive.Set(stats)
			}
		} else if decorator, ok := storage.(*mongodb.MetricsDecorator); ok {
			// Versuche den ursprünglichen Client aus dem Decorator zu bekommen
			if originalStorage := decorator.GetOriginalStorage(); originalStorage != nil {
				if mongoClient, ok := originalStorage.(*mongodb.Client); ok {
					if client := mongoClient.GetMongoClient(); client != nil {
						stats := float64(client.NumberSessionsInProgress())
						domains.DBConnectionsActive.Set(stats)
					}
				}
			}
		}
	}
}

func startPrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Msgf("Starting Prometheus metrics server on :%d", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Error().Err(err).Msg("Prometheus metrics server failed")
	}
}
