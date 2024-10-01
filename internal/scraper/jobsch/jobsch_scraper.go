package jobsch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"job-scraper/internal/models"
	"net/http"
	"time"
)

type JobschScraper struct {
	baseURL string
	apiKey  string
}

func NewJobschScraper(baseURL, apiKey string) *JobschScraper {
	return &JobschScraper{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (s *JobschScraper) Scrape(ctx context.Context) ([]models.Job, error) {
	url := fmt.Sprintf("%s?api_key=%s", s.baseURL, s.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobsResponse struct {
		Jobs []struct {
			ID          string    `json:"id"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Company     string    `json:"company"`
			Location    string    `json:"location"`
			PostedAt    time.Time `json:"posted_at"`
		} `json:"jobs"`
	}

	if err := json.Unmarshal(body, &jobsResponse); err != nil {
		return nil, err
	}

	var jobs []models.Job
	for _, j := range jobsResponse.Jobs {
		job := models.Job{
			URL:         fmt.Sprintf("https://www.jobs.ch/en/vacancies/%s", j.ID),
			Title:       j.Title,
			Description: j.Description,
			Company:     j.Company,
			Location:    j.Location,
			PostingDate: j.PostedAt,
			IsActive:    true,
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (s *JobschScraper) Name() string {
	return "Jobs.ch"
}
