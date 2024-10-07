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
	client    HTTPClient
	baseURL   string
	maxPages  int
	pageSize  int
	parseFunc func([]byte) (*models.Job, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	BaseURL   string
	MaxPages  int
	PageSize  int
	ParseFunc func([]byte) (*models.Job, error)
}

func NewJobsChScraper(config Config) *JobsChScraper {
	return &JobsChScraper{
		client:    &http.Client{Timeout: 10 * time.Second},
		baseURL:   config.BaseURL,
		maxPages:  config.MaxPages,
		pageSize:  config.PageSize,
		parseFunc: config.ParseFunc,
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

// TODO impl second api call for id
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
		job, err := s.parseFunc(doc)
		if err != nil {
			log.Warn().Err(err).Msg("Error parsing job")
			continue
		}
		jobs = append(jobs, *job)
	}

	return jobs, nil
}

func (s *JobsChScraper) Name() string {
	return "Jobs.ch"
}
