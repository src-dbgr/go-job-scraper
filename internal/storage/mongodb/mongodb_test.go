package mongodb

import (
	"context"
	"job-scraper/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockJobRepository struct {
	mock.Mock
}

func (m *MockJobRepository) SaveJob(ctx context.Context, job models.Job) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *MockJobRepository) GetJobByID(ctx context.Context, id string) (*models.Job, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Job), args.Error(1)
}

func (m *MockJobRepository) GetJobs(ctx context.Context) ([]models.Job, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Job), args.Error(1)
}

func (m *MockJobRepository) GetJobCountByCategory(ctx context.Context) (map[string]int, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockJobRepository) GetTotalJobCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockJobRepository) GetExistingURLs(ctx context.Context) (map[string]bool, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]bool), args.Error(1)
}

func (m *MockJobRepository) AggregateJobs(ctx context.Context, pipeline mongo.Pipeline) ([]bson.M, error) {
	args := m.Called(ctx, pipeline)
	return args.Get(0).([]bson.M), args.Error(1)
}

func TestSaveJob(t *testing.T) {
	mockRepo := new(MockJobRepository)
	job := models.Job{
		ID:          primitive.NewObjectID(),
		Title:       "Test Job",
		Company:     "Test Company",
		PostingDate: time.Now(),
	}

	mockRepo.On("SaveJob", mock.Anything, job).Return(nil)

	err := mockRepo.SaveJob(context.Background(), job)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetJobByID(t *testing.T) {
	mockRepo := new(MockJobRepository)
	job := &models.Job{
		ID:          primitive.NewObjectID(),
		Title:       "Test Job",
		Company:     "Test Company",
		PostingDate: time.Now(),
	}

	mockRepo.On("GetJobByID", mock.Anything, job.ID.Hex()).Return(job, nil)

	result, err := mockRepo.GetJobByID(context.Background(), job.ID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, job, result)
	mockRepo.AssertExpectations(t)
}
