package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"job-scraper/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockPromptRepository struct {
	mock.Mock
}

func (m *MockPromptRepository) GetPrompt(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func TestProcessor_Process(t *testing.T) {
	// Create a mock HTTP client
	mockClient := &http.Client{
		Transport: &mockTransport{},
	}

	// Create a mock prompt repository
	mockRepo := new(MockPromptRepository)

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": `{
							"title": "Test Job",
							"company": "Test Company",
							"postingDate": "2024-10-20T00:00:00Z",
							"expirationDate": "2024-11-20T00:00:00Z",
							"isActive": true,
							"jobCategories": ["SOFTWARE_DEVELOPER"],
							"mustSkills": ["Go", "MongoDB"],
							"optionalSkills": ["Docker", "Kubernetes"],
							"salary": "100000-120000",
							"yearsOfExperience": 3,
							"educationLevel": "Bachelor's",
							"benefits": ["Health Insurance", "401k"],
							"companySize": 500,
							"workCulture": "Agile",
							"remote": true,
							"languages": ["English", "Spanish"]
						}`,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	// Create a Processor instance with the mock client and repository
	processor := NewProcessor(Config{
		APIURL: ts.URL,
		APIKey: "test-key",
	}, mockRepo)
	processor.client = mockClient

	// Set up expectations
	mockRepo.On("GetPrompt", "job_extraction").Return("Test prompt", nil)

	// Call the Process method
	ctx := context.Background()
	job := models.Job{
		ID:          primitive.NewObjectID(),
		Description: "Test job description",
	}
	processedJob, err := processor.Process(ctx, job)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, "Test Job", processedJob.Title)
	assert.Equal(t, "Test Company", processedJob.Company)
	assert.Equal(t, time.Date(2024, 10, 20, 0, 0, 0, 0, time.UTC), processedJob.PostingDate)
	assert.Equal(t, time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC), processedJob.ExpirationDate)
	assert.True(t, processedJob.IsActive)
	assert.Equal(t, []string{"SOFTWARE_DEVELOPER"}, processedJob.JobCategories)
	assert.Equal(t, []string{"Go", "MongoDB"}, processedJob.MustSkills)
	assert.Equal(t, []string{"Docker", "Kubernetes"}, processedJob.OptionalSkills)
	assert.Equal(t, "100000-120000", processedJob.Salary)
	assert.Equal(t, 3, processedJob.YearsOfExperience)
	assert.Equal(t, "Bachelor's", processedJob.EducationLevel)
	assert.Equal(t, []string{"Health Insurance", "401k"}, processedJob.Benefits)
	assert.Equal(t, 500, processedJob.CompanySize)
	assert.Equal(t, "Agile", processedJob.WorkCulture)
	assert.True(t, processedJob.Remote)
	assert.Equal(t, []string{"English", "Spanish"}, processedJob.Languages)

	// Verify that the expectations were met
	mockRepo.AssertExpectations(t)
}

type mockTransport struct{}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response := map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]interface{}{
					"content": `{
						"title": "Test Job",
						"company": "Test Company",
						"postingDate": "2024-10-20T00:00:00Z",
						"expirationDate": "2024-11-20T00:00:00Z",
						"isActive": true,
						"jobCategories": ["SOFTWARE_DEVELOPER"],
						"mustSkills": ["Go", "MongoDB"],
						"optionalSkills": ["Docker", "Kubernetes"],
						"salary": "100000-120000",
						"yearsOfExperience": 3,
						"educationLevel": "Bachelor's",
						"benefits": ["Health Insurance", "401k"],
						"companySize": 500,
						"workCulture": "Agile",
						"remote": true,
						"languages": ["English", "Spanish"]
					}`,
				},
			},
		},
	}
	responseBody, _ := json.Marshal(response)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
		Header:     make(http.Header),
	}, nil
}
