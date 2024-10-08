package jobsch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"job-scraper/internal/models"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type JobsChScraper struct {
	client     HTTPClient
	baseURL    string
	maxPages   int
	pageSize   int
	parseFunc  func([]byte) (*models.Job, error)
	jobFetcher JobFetcher
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type JobFetcher interface {
	FetchJob(ctx context.Context, jobID string) (*models.Job, error)
}

type Config struct {
	BaseURL    string
	MaxPages   int
	PageSize   int
	ParseFunc  func([]byte) (*models.Job, error)
	JobFetcher JobFetcher
}

func NewJobsChScraper(config Config) *JobsChScraper {
	return &JobsChScraper{
		client:     &http.Client{Timeout: 10 * time.Second},
		baseURL:    config.BaseURL,
		maxPages:   config.MaxPages,
		pageSize:   config.PageSize,
		parseFunc:  config.ParseFunc,
		jobFetcher: config.JobFetcher,
	}
}

func (s *JobsChScraper) Scrape(ctx context.Context) ([]models.Job, error) {
	return s.ScrapePages(ctx, s.maxPages)
}

func (s *JobsChScraper) ScrapePages(ctx context.Context, pages int) ([]models.Job, error) {
	var allJobs []models.Job

	for page := 1; page <= pages; page++ {
		select {
		case <-ctx.Done():
			return allJobs, ctx.Err()
		default:
			jobs, err := s.scrapePage(ctx, page)
			if err != nil {
				log.Error().Err(err).Int("page", page).Msg("Error scraping page")
				continue
			}
			allJobs = append(allJobs, jobs...)
			if len(jobs) < s.pageSize {
				return allJobs, nil // No more jobs to scrape
			}
		}
	}

	return allJobs, nil
}

func (s *JobsChScraper) scrapePage(ctx context.Context, page int) ([]models.Job, error) {
	url := fmt.Sprintf("%s/public/search?page=%d&query=software&rows=%d", s.baseURL, page, s.pageSize)
	log.Info().Int("page", page).Msg("Getting JobsCh Page result")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var searchResponse struct {
		Documents []json.RawMessage `json:"documents"`
	}
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	var jobs []models.Job
	for _, doc := range searchResponse.Documents {
		var jobData struct {
			JobID string `json:"job_id"`
		}
		if err := json.Unmarshal(doc, &jobData); err != nil {
			log.Warn().Err(err).Msg("Error parsing job data")
			continue
		}

		job, err := s.jobFetcher.FetchJob(ctx, jobData.JobID)
		if err != nil {
			log.Warn().Err(err).Str("jobID", jobData.JobID).Msg("Error fetching job details")
			continue
		}

		jobs = append(jobs, *job)
	}

	return jobs, nil
}

func (s *JobsChScraper) Name() string {
	return "Jobs.ch"
}

// JobsChFetcher impls the JobFetcher Interface
type JobsChFetcher struct {
	client  HTTPClient
	baseURL string
}

func NewJobsChFetcher(client HTTPClient, baseURL string) *JobsChFetcher {
	return &JobsChFetcher{
		client:  client,
		baseURL: baseURL,
	}
}

func (f *JobsChFetcher) FetchJob(ctx context.Context, jobID string) (*models.Job, error) {
	url := fmt.Sprintf("%s/public/search/job/%s", f.baseURL, jobID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	job := &models.Job{
		URL:         url,
		Description: string(body),
	}

	return job, nil
}
