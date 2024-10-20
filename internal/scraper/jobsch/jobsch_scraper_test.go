package jobsch

import (
	"bytes"
	"context"
	"io"
	"job-scraper/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type MockJobFetcher struct {
	mock.Mock
}

func (m *MockJobFetcher) FetchJob(ctx context.Context, jobID string) (*models.Job, error) {
	args := m.Called(ctx, jobID)
	return args.Get(0).(*models.Job), args.Error(1)
}

func TestJobsChScraper_Scrape(t *testing.T) {
	// Create a mock HTTP client
	mockClient := new(MockHTTPClient)
	mockFetcher := new(MockJobFetcher)

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"documents": [{"job_id": "123"}]}`))
	}))
	defer ts.Close()

	// Create a JobsChScraper instance with the mock client
	scraper := NewJobsChScraper(Config{
		BaseURL:    ts.URL,
		MaxPages:   1,
		PageSize:   10,
		JobFetcher: mockFetcher,
	})
	scraper.client = mockClient

	// Set up expectations
	mockClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"documents": [{"job_id": "123"}]}`))),
	}, nil)

	objectID := primitive.NewObjectID()
	mockFetcher.On("FetchJob", mock.Anything, "123").Return(&models.Job{
		ID:    objectID,
		Title: "Test Job",
	}, nil)

	// Call the Scrape method
	ctx := context.Background()
	jobs, err := scraper.Scrape(ctx)

	// Assert the results
	assert.NoError(t, err)
	assert.Len(t, jobs, 1)
	assert.Equal(t, objectID, jobs[0].ID)
	assert.Equal(t, "Test Job", jobs[0].Title)

	// Verify that the expectations were met
	mockClient.AssertExpectations(t)
	mockFetcher.AssertExpectations(t)
}
