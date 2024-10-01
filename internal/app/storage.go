package app

import (
	"context"
	"job-scraper/internal/config"
	"job-scraper/internal/storage"
	"job-scraper/internal/storage/mongodb"
)

func initStorage(ctx context.Context, cfg *config.Config) (storage.Storage, error) {
	return mongodb.NewClient(ctx, cfg.MongoDB.URI, cfg.MongoDB.Database)
}
