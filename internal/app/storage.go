package app

import (
	"context"
	"job-scraper/internal/config"
	"job-scraper/internal/storage"
	"job-scraper/internal/storage/mongodb"
)

func initStorage(ctx context.Context, cfg *config.Config) (storage.Storage, error) {
	baseStorage, err := mongodb.NewClient(ctx, cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		return nil, err
	}

	// wrape the base storage to the metricsdecorator
	return mongodb.NewMetricsDecorator(baseStorage), nil
}
