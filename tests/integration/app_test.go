package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"job-scraper/internal/app"
	"job-scraper/internal/models"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// getFreePort findet einen freien Port
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

type TestEnv struct {
	mongoContainer testcontainers.Container
	mockServer     *http.Server
	apiPort        int
	mockPort       int
}

func setupMockServer(port int) *http.Server {
	mux := http.NewServeMux()

	// Mock für die Suche mit verschiedenen Jobs
	mux.HandleFunc("/api/v1/public/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"documents": []map[string]interface{}{
				{
					"job_id":   "dev-1",
					"title":    "Software Engineer",
					"company":  "Tech Corp",
					"location": "Zürich",
				},
				{
					"job_id":   "dev-2",
					"title":    "DevOps Engineer",
					"company":  "Cloud Systems AG",
					"location": "Bern",
				},
				{
					"job_id":   "data-1",
					"title":    "Data Scientist",
					"company":  "AI Solutions GmbH",
					"location": "Basel",
				},
				// Wiederhole die Job-Muster, um auf 20 zu kommen
			},
		})
	})

	// Mock für Job-Details basierend auf Job-ID Präfix
	mux.HandleFunc("/api/v1/public/search/job/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jobID := filepath.Base(r.URL.Path)

		var jobDetails map[string]interface{}

		// Bestimme Job-Typ basierend auf Job-ID
		switch {
		case strings.HasPrefix(jobID, "dev-1"):
			jobDetails = map[string]interface{}{
				"title":          "Software Engineer",
				"company":        "Tech Corp",
				"description":    "Looking for a senior engineer with strong backend experience...",
				"location":       "Zürich",
				"employmentType": "Full-time",
				"salary":         "120000-160000 CHF",
			}
		case strings.HasPrefix(jobID, "dev-2"):
			jobDetails = map[string]interface{}{
				"title":          "DevOps Engineer",
				"company":        "Cloud Systems AG",
				"description":    "DevOps position with focus on Kubernetes...",
				"location":       "Bern",
				"employmentType": "Full-time",
				"salary":         "110000-150000 CHF",
			}
		default:
			jobDetails = map[string]interface{}{
				"title":          "Data Scientist",
				"company":        "AI Solutions GmbH",
				"description":    "Data Science position with focus on ML...",
				"location":       "Basel",
				"employmentType": "80-100%",
				"salary":         "115000-155000 CHF",
			}
		}

		json.NewEncoder(w).Encode(jobDetails)
	})

	// Mock für OpenAI mit job-spezifischen Antworten
	mux.HandleFunc("/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var requestBody struct {
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		json.NewDecoder(r.Body).Decode(&requestBody)

		var response string
		if strings.Contains(requestBody.Messages[0].Content, "Software Engineer") {
			response = `{
                "title": "Software Engineer",
                "company": "Tech Corp",
                "description": "Backend engineering position",
                "location": "Zürich",
                "employmentType": "Full-time",
                "postingDate": "2024-10-22T00:00:00Z",
                "expirationDate": "2024-11-22T00:00:00Z",
                "isActive": true,
                "jobCategories": ["SOFTWARE_DEVELOPER"],
                "mustSkills": ["Go", "Java", "Microservices"],
                "optionalSkills": ["Kubernetes", "AWS"],
                "salary": "120000-160000 CHF",
                "yearsOfExperience": 5,
                "educationLevel": "Master's",
                "benefits": ["Health Insurance", "Remote Work"],
                "companySize": 500,
                "workCulture": "Agile",
                "remote": true,
                "languages": ["English", "German"]
            }`
		} else if strings.Contains(requestBody.Messages[0].Content, "DevOps") {
			response = `{
                "title": "DevOps Engineer",
                "company": "Cloud Systems AG",
                "description": "DevOps position",
                "location": "Bern",
                "employmentType": "Full-time",
                "postingDate": "2024-10-22T00:00:00Z",
                "expirationDate": "2024-11-22T00:00:00Z",
                "isActive": true,
                "jobCategories": ["DEVOPS_ENGINEER"],
                "mustSkills": ["Kubernetes", "Docker", "Jenkins"],
                "optionalSkills": ["Terraform", "AWS"],
                "salary": "110000-150000 CHF",
                "yearsOfExperience": 3,
                "educationLevel": "Bachelor's",
                "benefits": ["Flexible Hours"],
                "companySize": 200,
                "workCulture": "DevOps",
                "remote": true,
                "languages": ["English", "German"]
            }`
		} else {
			response = `{
                "title": "Data Scientist",
                "company": "AI Solutions GmbH",
                "description": "Data Science position",
                "location": "Basel",
                "employmentType": "80-100%",
                "postingDate": "2024-10-22T00:00:00Z",
                "expirationDate": "2024-11-22T00:00:00Z",
                "isActive": true,
                "jobCategories": ["DATA_SCIENTIST"],
                "mustSkills": ["Python", "Machine Learning", "Deep Learning"],
                "optionalSkills": ["TensorFlow", "PyTorch"],
                "salary": "115000-155000 CHF",
                "yearsOfExperience": 4,
                "educationLevel": "PhD",
                "benefits": ["Research Budget"],
                "companySize": 100,
                "workCulture": "Research-Driven",
                "remote": false,
                "languages": ["English", "German"]
            }`
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]interface{}{
						"content": response,
					},
				},
			},
		})
	})

	// Server starten...
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Mock server error: %v\n", err)
		}
	}()

	// Warten bis Server bereit ist...
	ready := make(chan bool)
	go func() {
		for i := 0; i < 50; i++ {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/v1/public/search", port))
			if err == nil {
				resp.Body.Close()
				ready <- true
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
		ready <- false
	}()

	if !<-ready {
		panic("Mock server failed to start")
	}

	return server
}

// In setupTestEnv auch Logging hinzufügen:
func setupTestEnv(t *testing.T) (*TestEnv, error) {
	t.Log("Setting up mock server...")
	mockPort, err := getFreePort()
	require.NoError(t, err, "Failed to get free port for mock server")
	t.Logf("Using mock server port: %d", mockPort)

	t.Log("Setting up API server...")
	apiPort, err := getFreePort()
	require.NoError(t, err, "Failed to get free port for API server")
	t.Logf("Using API server port: %d", apiPort)

	t.Log("Starting MongoDB container...")
	mongoContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:6.0",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
		},
		Started: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	mappedPort, _ := mongoContainer.MappedPort(context.Background(), "27017")
	hostIP, _ := mongoContainer.Host(context.Background())
	mongoURI := fmt.Sprintf("mongodb://%s:%s", hostIP, mappedPort.Port())

	mockServer := setupMockServer(mockPort)
	mockBaseURL := fmt.Sprintf("http://localhost:%d", mockPort)

	// Set environment variables
	t.Setenv("MONGODB_URI", mongoURI)
	t.Setenv("MONGODB_DATABASE", "jobsdb_test")
	t.Setenv("API_PORT", fmt.Sprintf("%d", apiPort))
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("OPENAI_API_URL", fmt.Sprintf("%s/chat/completions", mockBaseURL))
	t.Setenv("SCRAPER_JOBSCH_BASE_URL", fmt.Sprintf("%s/api/v1", mockBaseURL))
	t.Setenv("PROMETHEUS_PORT", "0")

	t.Log("Test environment setup completed")
	return &TestEnv{
		mongoContainer: mongoContainer,
		mockServer:     mockServer,
		mockPort:       mockPort,
		apiPort:        apiPort,
	}, nil
}

func (env *TestEnv) Cleanup(ctx context.Context) error {
	if env.mockServer != nil {
		env.mockServer.Shutdown(ctx)
	}
	if env.mongoContainer != nil {
		return env.mongoContainer.Terminate(ctx)
	}
	return nil
}

// Hilfsfunktion für Job-Validierung
func validateJobProperties(t *testing.T, job models.Job) {
	t.Helper()

	// Grundlegende Job-Eigenschaften sollten vorhanden sein
	assert.NotEmpty(t, job.Title, "Job title should not be empty")
	assert.NotEmpty(t, job.Company, "Company should not be empty")
	assert.NotEmpty(t, job.Location, "Location should not be empty")
	assert.NotEmpty(t, job.JobCategories, "Job categories should not be empty")
	assert.NotEmpty(t, job.MustSkills, "Must skills should not be empty")
	assert.NotEmpty(t, job.OptionalSkills, "Optional skills should not be empty")
	assert.NotEmpty(t, job.EmploymentType, "Employment type should not be empty")

	// Zeitstempel sollten gültig sein
	assert.False(t, job.PostingDate.IsZero(), "Posting date should be set")
	assert.False(t, job.ExpirationDate.IsZero(), "Expiration date should be set")
	assert.True(t, job.IsActive, "Job should be active")

	// Validiere Job-Kategorie
	assert.Contains(t, []string{
		"SOFTWARE_DEVELOPER",
		"DEVOPS_ENGINEER",
		"DATA_SCIENTIST",
	}, job.JobCategories[0], "Job category should be valid")
}

func TestJobScraper(t *testing.T) {
	t.Log("Starting integration test...")

	// Get project root and set config path
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	require.NoError(t, err, "Failed to get project root path")

	configPath := filepath.Join(projectRoot, "configs")
	t.Logf("Using config path: %s", configPath)

	_, err = os.Stat(filepath.Join(configPath, "config.yaml"))
	require.NoError(t, err, "Config file should exist at %s", configPath)

	t.Setenv("JOBSCRAPER_CONFIG_PATH", configPath)

	// Setup test environment
	t.Log("Setting up test environment...")
	ctx := context.Background()
	env, err := setupTestEnv(t)
	require.NoError(t, err, "Failed to setup test environment")

	defer func() {
		t.Log("Cleaning up test environment...")
		err := env.Cleanup(ctx)
		if err != nil {
			t.Logf("Failed to cleanup test environment: %v", err)
		}
	}()

	// Log environment details
	t.Logf("Test environment configured with:")
	t.Logf("- Mock server port: %d", env.mockPort)
	t.Logf("- API server port: %d", env.apiPort)
	t.Logf("- MongoDB URI: %s", os.Getenv("MONGODB_URI"))

	// Initialize application
	t.Log("Initializing application...")
	application, err := app.New(ctx)
	require.NoError(t, err, "Failed to initialize application")

	// Run application
	t.Log("Starting application...")
	go application.Run(ctx)
	defer func() {
		t.Log("Shutting down application...")
		application.Shutdown(ctx)
	}()

	// Wait for application to start and be ready
	baseURL := fmt.Sprintf("http://localhost:%d", env.apiPort)
	t.Logf("Waiting for application to be ready at %s", baseURL)

	var serverReady bool
	for i := 0; i < 50; i++ {
		if _, err := http.Get(baseURL + "/api/v1/jobs"); err == nil {
			serverReady = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	require.True(t, serverReady, "Server should be ready")
	t.Log("Application is ready")

	// Run the scraping test
	t.Run("Scraping Process", func(t *testing.T) {
		client := &http.Client{Timeout: 10 * time.Second}

		// Trigger scraping
		t.Log("Triggering job scraping...")
		scrapeURL := fmt.Sprintf("%s/api/v1/scrape/jobsch?pages=1", baseURL)
		t.Logf("POST %s", scrapeURL)

		req, err := http.NewRequestWithContext(ctx, "POST", scrapeURL, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusAccepted, resp.StatusCode, "Scraping request should be accepted")
		resp.Body.Close()

		// Wait for scraping to complete
		t.Log("Waiting for scraping to complete...")
		time.Sleep(5 * time.Second)

		// Verify jobs were scraped
		t.Log("Verifying scraped jobs...")
		req, err = http.NewRequestWithContext(ctx, "GET",
			fmt.Sprintf("%s/api/v1/jobs", baseURL), nil)
		require.NoError(t, err)

		resp, err = client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, "Should get OK status when fetching jobs")

		var jobs []models.Job
		err = json.NewDecoder(resp.Body).Decode(&jobs)
		require.NoError(t, err, "Should be able to decode jobs response")
		resp.Body.Close()

		// Prüfe die Anzahl der Jobs (eine Seite = 20 Jobs)
		assert.Equal(t, 20, len(jobs), "Should have found exactly 20 jobs (one page)")

		// Log job details
		t.Logf("Found %d jobs", len(jobs))
		for i, job := range jobs {
			t.Logf("Job %d:", i+1)
			t.Logf("  - Title: %s", job.Title)
			t.Logf("  - Company: %s", job.Company)
			t.Logf("  - Location: %s", job.Location)
			t.Logf("  - Categories: %v", job.JobCategories)
			t.Logf("  - Must Skills: %v", job.MustSkills)
			t.Logf("  - Optional Skills: %v", job.OptionalSkills)

			// Validiere jeden Job
			validateJobProperties(t, job)
		}
	})

	// Test job statistics
	t.Run("Job Statistics", func(t *testing.T) {
		t.Log("Testing job statistics...")
		client := &http.Client{Timeout: 10 * time.Second}

		statsURL := fmt.Sprintf("%s/api/v1/stats/job-categories-counts", baseURL)
		t.Logf("GET %s", statsURL)

		req, err := http.NewRequestWithContext(ctx, "GET", statsURL, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, "Should get OK status when fetching stats")

		var stats []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&stats)
		require.NoError(t, err, "Should be able to decode stats response")
		resp.Body.Close()

		t.Logf("Found statistics: %+v", stats)
		assert.NotEmpty(t, stats, "Should have job statistics")
	})

	t.Log("Integration test completed successfully")
}
