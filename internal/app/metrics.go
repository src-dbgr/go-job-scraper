package app

import (
	"fmt"
	"job-scraper/internal/metrics"
	"job-scraper/internal/metrics/domains"
	"job-scraper/internal/storage"
	"job-scraper/internal/storage/mongodb"
	"net"
	"net/http"
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

// startPrometheusServer initiates the metrics server
func startPrometheusServer(port int) error {
	server, err := setupServer(port)
	if err != nil {
		return err
	}

	// Start server and verify it's running
	if err := runServerWithVerification(server, port); err != nil {
		return err
	}

	return nil
}

func setupServer(port int) (*http.Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create listener for Prometheus metrics server")
		return nil, err
	}

	server := &http.Server{
		Handler: promhttp.Handler(),
	}

	log.Info().Msgf("Starting Prometheus metrics server on port %d", port)

	// Start server in go routine
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Prometheus metrics server failed")
		}
	}()

	return server, nil
}

func runServerWithVerification(server *http.Server, port int) error {
	// Give the server a moment to start
	time.Sleep(150 * time.Millisecond)

	// Verify server is responding
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/metrics", port))
	if err != nil {
		log.Error().Err(err).Msg("Failed to verify Prometheus metrics server")
		return err
	}
	defer resp.Body.Close()

	log.Info().Msgf("Started Prometheus metrics server on port %d", port)
	return nil
}
