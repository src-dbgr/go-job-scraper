package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"job-scraper/internal/apperrors"
	"job-scraper/internal/models"
	"job-scraper/internal/parser"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Processor struct {
	client     HTTPClient
	config     Config
	promptRepo PromptRepository
	jobParser  *parser.JobParser
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	APIURL      string
	APIKey      string
	Model       string
	Timeout     time.Duration
	Temperature float64
	MaxTokens   int
	TopP        float64
	FreqPenalty float64
	PresPenalty float64
}

type PromptRepository interface {
	GetPrompt(name string) (string, error)
}

func NewProcessor(config Config, promptRepo PromptRepository) *Processor {
	return &Processor{
		client:     &http.Client{},
		config:     config,
		promptRepo: promptRepo,
		jobParser:  parser.NewJobParser(),
	}
}

func (p *Processor) Process(ctx context.Context, job models.Job) (models.Job, error) {
	updatedJob, err := p.extractJobInfo(ctx, job.Description)
	if err != nil {
		return job, apperrors.NewProcessingError(
			job.ID.Hex(),
			"failed to process job with OpenAI",
			err,
		)
	}

	// Preserve the original URL and any other fields that should not be overwritten
	updatedJob.URL = job.URL

	log.Info().
		Str("job_title", updatedJob.Title).
		Strs("extracted_skills", updatedJob.MustSkills).
		Msg("Processed job with OpenAI")

	return *updatedJob, nil
}

func (p *Processor) extractJobInfo(ctx context.Context, jobDescription string) (*models.Job, error) {
	prompt, err := p.promptRepo.GetPrompt("job_extraction")
	if err != nil {
		return nil, apperrors.NewProcessingError("", "Failed to get prompt", err)
	}

	payload := map[string]interface{}{
		"model": p.config.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf(prompt, jobDescription),
			},
		},
		"temperature":       p.config.Temperature,
		"max_tokens":        p.config.MaxTokens,
		"top_p":             p.config.TopP,
		"frequency_penalty": p.config.FreqPenalty,
		"presence_penalty":  p.config.PresPenalty,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, apperrors.NewProcessingError("", "Failed marshalling payload", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.config.APIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, apperrors.NewProcessingError("", "Failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, apperrors.NewProcessingError("", "Failed to make request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, apperrors.NewProcessingError(
			"",
			fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
			fmt.Errorf("body: %s", string(body)),
		)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, apperrors.NewProcessingError("", "Error decoding response", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	jsonContent, err := extractJSONFromContent(result.Choices[0].Message.Content)
	if err != nil {
		return nil, apperrors.NewProcessingError("", "Failed to extract JSON from OpenAI response", err)
	}

	job, err := p.jobParser.ParseJob([]byte(jsonContent))
	if err != nil {
		return nil, apperrors.NewProcessingError("", "Failed to parse job information", err)
	}

	return job, nil
}

func extractJSONFromContent(content string) (string, error) {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	return strings.TrimSpace(content), nil
}
