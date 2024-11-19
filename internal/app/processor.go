package app

import (
	"fmt"
	"job-scraper/internal/config"
	"job-scraper/internal/processor"
	"job-scraper/internal/processor/openai"
	// Zukünftige Prozessor-Implementierungen:
	// "job-scraper/internal/processor/claude"
	// "job-scraper/internal/processor/gpt4all"
	// etc.
)

func initProcessor(cfg *config.Config) (processor.JobProcessor, error) {
	switch cfg.Processor.Type {
	case "openai":
		return initOpenAIProcessor(cfg)
	// Zukünftige Prozessor-Typen:
	// case "claude":
	//     return initClaudeProcessor(cfg)
	// case "gemini":
	//     return initGeminiProcessor(cfg)
	default:
		return nil, fmt.Errorf("unsupported processor type: %s", cfg.Processor.Type)
	}
}

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
