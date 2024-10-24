package mongodb

import (
	"context"
	"time"

	"job-scraper/internal/metrics/domains"
	"job-scraper/internal/models"
	"job-scraper/internal/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MetricsDecorator struct {
	storage storage.Storage
}

func NewMetricsDecorator(storage storage.Storage) storage.Storage {
	return &MetricsDecorator{storage: storage}
}

func (d *MetricsDecorator) GetOriginalStorage() storage.Storage {
	return d.storage
}

func (d *MetricsDecorator) SaveJob(ctx context.Context, job models.Job) error {
	start := time.Now()
	err := d.storage.SaveJob(ctx, job)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	domains.DBOperationDuration.WithLabelValues("save_job", status).Observe(duration)
	domains.DBOperationsTotal.WithLabelValues("save_job", status).Inc()

	return err
}

func (d *MetricsDecorator) GetJobs(ctx context.Context) ([]models.Job, error) {
	start := time.Now()
	jobs, err := d.storage.GetJobs(ctx)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	domains.DBOperationDuration.WithLabelValues("get_jobs", status).Observe(duration)
	domains.DBOperationsTotal.WithLabelValues("get_jobs", status).Inc()

	return jobs, err
}

func (d *MetricsDecorator) GetJobByID(ctx context.Context, id string) (*models.Job, error) {
	start := time.Now()
	job, err := d.storage.GetJobByID(ctx, id)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	domains.DBOperationDuration.WithLabelValues("get_job_by_id", status).Observe(duration)
	domains.DBOperationsTotal.WithLabelValues("get_job_by_id", status).Inc()

	return job, err
}

func (d *MetricsDecorator) GetJobCountByCategory(ctx context.Context) (map[string]int, error) {
	start := time.Now()
	counts, err := d.storage.GetJobCountByCategory(ctx)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	domains.DBOperationDuration.WithLabelValues("get_job_count_by_category", status).Observe(duration)
	domains.DBOperationsTotal.WithLabelValues("get_job_count_by_category", status).Inc()

	return counts, err
}

func (d *MetricsDecorator) GetTotalJobCount(ctx context.Context) (int, error) {
	start := time.Now()
	count, err := d.storage.GetTotalJobCount(ctx)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	domains.DBOperationDuration.WithLabelValues("get_total_job_count", status).Observe(duration)
	domains.DBOperationsTotal.WithLabelValues("get_total_job_count", status).Inc()

	return count, err
}

func (d *MetricsDecorator) GetExistingURLs(ctx context.Context) (map[string]bool, error) {
	start := time.Now()
	urls, err := d.storage.GetExistingURLs(ctx)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	domains.DBOperationDuration.WithLabelValues("get_existing_urls", status).Observe(duration)
	domains.DBOperationsTotal.WithLabelValues("get_existing_urls", status).Inc()

	return urls, err
}

func (d *MetricsDecorator) Close(ctx context.Context) error {
	start := time.Now()
	err := d.storage.Close(ctx)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	domains.DBOperationDuration.WithLabelValues("close", status).Observe(duration)
	domains.DBOperationsTotal.WithLabelValues("close", status).Inc()

	return err
}

func (d *MetricsDecorator) AggregateJobs(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error) {
	start := time.Now()
	results, err := d.storage.AggregateJobs(ctx, pipeline)
	duration := time.Since(start).Seconds()

	status := "success"
	if err != nil {
		status = "error"
	}

	domains.DBOperationDuration.WithLabelValues("aggregate_jobs", status).Observe(duration)
	domains.DBOperationsTotal.WithLabelValues("aggregate_jobs", status).Inc()

	return results, err
}
