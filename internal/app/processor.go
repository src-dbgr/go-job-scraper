package app

import (
	"fmt"
	"job-scraper/internal/config"
	"job-scraper/internal/processor"
	"job-scraper/internal/processor/openai"
	// Future processor implementations:
	// "job-scraper/internal/processor/claude"
	// "job-scraper/internal/processor/gpt4all"
	// etc.
)

// initProcessor initializes the appropriate job processor based on the configuration
// Returns a JobProcessor interface implementation and an error if initialization fails
func initProcessor(cfg *config.Config) (processor.JobProcessor, error) {
	switch cfg.Processor.Type {
	case "openai":
		return initOpenAIProcessor(cfg)
	// Future processor types:
	// case "claude":
	//     return initClaudeProcessor(cfg)
	// case "gemini":
	//     return initGeminiProcessor(cfg)
	default:
		return nil, fmt.Errorf("unsupported processor type: %s", cfg.Processor.Type)
	}
}

// initOpenAIProcessor initializes an OpenAI processor with the provided configuration
// Returns a configured OpenAI processor instance and an error if initialization fails
func initOpenAIProcessor(cfg *config.Config) (processor.JobProcessor, error) {
	openaiConfig := openai.Config{
		APIURL:      cfg.OpenAI.APIURL,
		APIKey:      cfg.OpenAI.APIKey,
		Model:       cfg.OpenAI.Model,
		Temperature: cfg.OpenAI.Temperature,
		MaxTokens:   cfg.OpenAI.MaxTokens,
		TopP:        cfg.OpenAI.TopP,
		FreqPenalty: cfg.OpenAI.FreqPenalty,
		PresPenalty: cfg.OpenAI.PresPenalty,
	}
	promptRepo := openai.NewFilePromptRepository()
	return openai.NewProcessor(openaiConfig, promptRepo), nil
}
