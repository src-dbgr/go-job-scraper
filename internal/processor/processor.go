package processor

import (
	"context"
	"job-scraper/internal/models"
)

type JobProcessor interface {
	Process(ctx context.Context, job models.Job) (models.Job, error)
}
