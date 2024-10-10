package storage

import (
	"context"
	"job-scraper/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	GetJobs(ctx context.Context) ([]models.Job, error)
	GetJobByID(ctx context.Context, id string) (*models.Job, error)
	SaveJob(ctx context.Context, job models.Job) error
	GetJobCountByCategory(ctx context.Context) (map[string]int, error)
	GetTotalJobCount(ctx context.Context) (int, error)
	GetExistingURLs(ctx context.Context) (map[string]bool, error)
	Close(ctx context.Context) error
	AggregateJobs(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error)
}
