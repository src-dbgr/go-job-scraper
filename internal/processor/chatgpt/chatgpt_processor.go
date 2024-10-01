package chatgpt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"job-scraper/internal/models"
	"job-scraper/internal/processor"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type ChatGPTProcessor struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

type ChatGPTConfig struct {
	APIKey string
	APIURL string
}

func NewChatGPTProcessor(config ChatGPTConfig) processor.JobProcessor {
	return &ChatGPTProcessor{
		apiKey: config.APIKey,
		apiURL: config.APIURL,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (p *ChatGPTProcessor) Process(ctx context.Context, job models.Job) (models.Job, error) {
	prompt := p.buildPrompt(job)
	response, err := p.sendRequest(ctx, prompt)
	if err != nil {
		return job, fmt.Errorf("failed to process job with ChatGPT: %w", err)
	}

	return p.updateJobWithResponse(job, response)
}

func (p *ChatGPTProcessor) buildPrompt(job models.Job) string {
	return fmt.Sprintf("Analyze the following job description and extract key information:\n\nTitle: %s\nDescription: %s", job.Title, job.Description)
}

func (p *ChatGPTProcessor) sendRequest(ctx context.Context, prompt string) (string, error) {
	requestBody := fmt.Sprintf(`{
		"model": "gpt-3.5-turbo",
		"messages": [{"role": "user", "content": "%s"}]
	}`, prompt)

	req, err := http.NewRequestWithContext(ctx, "POST", p.apiURL, strings.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return result.Choices[0].Message.Content, nil
}

func (p *ChatGPTProcessor) updateJobWithResponse(job models.Job, response string) (models.Job, error) {

	// example
	if strings.Contains(response, "Skills:") {
		skills := strings.Split(strings.Split(response, "Skills:")[1], ",")
		job.MustSkills = make([]string, 0, len(skills))
		for _, skill := range skills {
			job.MustSkills = append(job.MustSkills, strings.TrimSpace(skill))
		}
	}

	// Log the processing result
	log.Info().
		Str("job_title", job.Title).
		Strs("extracted_skills", job.MustSkills).
		Msg("Processed job with ChatGPT")

	return job, nil
}
