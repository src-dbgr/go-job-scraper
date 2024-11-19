package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

type ScraperConfig struct {
	BaseURL      string `mapstructure:"base_url"`
	APIKey       string `mapstructure:"api_key"`
	Schedule     string `mapstructure:"schedule"`
	DefaultPages int    `mapstructure:"default_pages"`
	MaxPages     int    `mapstructure:"max_pages"`
}

type Config struct {
	API struct {
		Port int
	}
	MongoDB struct {
		URI      string
		Database string
	}
	Scrapers  map[string]*ScraperConfig
	Processor struct {
		Type string // "openai", "claude", "gpt4all", etc.
	}
	OpenAI struct {
		APIKey      string
		APIURL      string
		Model       string
		Timeout     time.Duration
		Temperature float64
		MaxTokens   int
		TopP        float64
		FreqPenalty float64
		PresPenalty float64
	}
	Logging struct {
		Level string
		File  string
	}
	Prometheus struct {
		Port int
	}
}

func LoadConfig() (*Config, error) {
	if configPath := os.Getenv("JOBSCRAPER_CONFIG_PATH"); configPath != "" {
		viper.AddConfigPath(configPath)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")

	// Environment-Variable Setup
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Lese zuerst die Konfigurationsdatei
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}

	// MongoDB configuration
	config.MongoDB.URI = viper.GetString("mongodb.uri")
	config.MongoDB.Database = viper.GetString("mongodb.database")

	// Validiere MongoDB URI
	if !strings.HasPrefix(config.MongoDB.URI, "mongodb://") &&
		!strings.HasPrefix(config.MongoDB.URI, "mongodb+srv://") {
		return nil, fmt.Errorf("invalid mongodb URI format")
	}

	// API configuration
	config.API.Port = viper.GetInt("api.port")
	if config.API.Port == 0 {
		config.API.Port = 8080
	}

	// Processor configuration
	config.Processor.Type = viper.GetString("processor.type")
	if config.Processor.Type == "" {
		config.Processor.Type = "openai"
	}

	// OpenAI configuration
	config.OpenAI.APIKey = viper.GetString("openai.api_key")
	config.OpenAI.APIURL = viper.GetString("openai.api_url")
	config.OpenAI.Model = viper.GetString("openai.model")
	config.OpenAI.Timeout = viper.GetDuration("openai.timeout")
	config.OpenAI.Temperature = viper.GetFloat64("openai.temperature")
	config.OpenAI.MaxTokens = viper.GetInt("openai.max_tokens")
	config.OpenAI.TopP = viper.GetFloat64("openai.top_p")
	config.OpenAI.FreqPenalty = viper.GetFloat64("openai.frequency_penalty")
	config.OpenAI.PresPenalty = viper.GetFloat64("openai.presence_penalty")

	// Logging configuration
	config.Logging.Level = viper.GetString("logging.level")
	config.Logging.File = viper.GetString("logging.file")

	// Prometheus configuration
	config.Prometheus.Port = viper.GetInt("prometheus.port")
	if config.Prometheus.Port == 0 {
		config.Prometheus.Port = 2112
	}

	// Scrapers configuration
	config.Scrapers = make(map[string]*ScraperConfig)
	scraperConfigs := viper.GetStringMap("scrapers")
	for scraperName, sc := range scraperConfigs {
		scraperConfigMap, ok := sc.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid scraper configuration for %s", scraperName)
		}

		cfg := &ScraperConfig{
			BaseURL:      getConfigValue(fmt.Sprintf("scrapers.%s.base_url", scraperName), scraperConfigMap["base_url"]),
			APIKey:       getConfigValue(fmt.Sprintf("scrapers.%s.api_key", scraperName), scraperConfigMap["api_key"]),
			Schedule:     viper.GetString(fmt.Sprintf("scrapers.%s.schedule", scraperName)),
			DefaultPages: viper.GetInt(fmt.Sprintf("scrapers.%s.default_pages", scraperName)),
			MaxPages:     viper.GetInt(fmt.Sprintf("scrapers.%s.max_pages", scraperName)),
		}

		// Validiere required fields
		if cfg.BaseURL == "" {
			return nil, &RequiredConfigError{Field: fmt.Sprintf("scrapers.%s.base_url", scraperName)}
		}
		if cfg.APIKey == "" {
			return nil, &RequiredConfigError{Field: fmt.Sprintf("scrapers.%s.api_key", scraperName)}
		}
		if cfg.Schedule == "" {
			return nil, &RequiredConfigError{Field: fmt.Sprintf("scrapers.%s.schedule", scraperName)}
		}

		// Default values für nicht-required fields
		if cfg.DefaultPages <= 0 {
			cfg.DefaultPages = 5
		}
		if cfg.MaxPages <= 0 {
			cfg.MaxPages = 20
		}

		// Validiere Schedule Format
		if _, err := cron.ParseStandard(cfg.Schedule); err != nil {
			return nil, fmt.Errorf("invalid schedule format for scraper %s: %w", scraperName, err)
		}

		config.Scrapers[scraperName] = cfg
	}

	return config, nil
}

func loadScraperConfig(config *Config) error {
	config.Scrapers = make(map[string]*ScraperConfig)

	scraperConfigs := viper.GetStringMap("scrapers")
	for scraperName, sc := range scraperConfigs {
		scraperConfigMap, ok := sc.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid scraper configuration for %s", scraperName)
		}

		// Required fields für jeden Scraper prüfen
		required := []string{"base_url", "api_key", "schedule"}
		for _, field := range required {
			if _, exists := scraperConfigMap[field]; !exists {
				return &RequiredConfigError{
					Field: fmt.Sprintf("scrapers.%s.%s", scraperName, field),
				}
			}
		}

		cfg := &ScraperConfig{
			BaseURL:      viper.GetString(fmt.Sprintf("scrapers.%s.base_url", scraperName)),
			APIKey:       viper.GetString(fmt.Sprintf("scrapers.%s.api_key", scraperName)),
			Schedule:     viper.GetString(fmt.Sprintf("scrapers.%s.schedule", scraperName)),
			DefaultPages: viper.GetInt(fmt.Sprintf("scrapers.%s.default_pages", scraperName)),
			MaxPages:     viper.GetInt(fmt.Sprintf("scrapers.%s.max_pages", scraperName)),
		}

		config.Scrapers[scraperName] = cfg
	}

	return nil
}

func validateConfig(config *Config) error {
	// Validiere MongoDB URI Format
	if !strings.HasPrefix(config.MongoDB.URI, "mongodb://") &&
		!strings.HasPrefix(config.MongoDB.URI, "mongodb+srv://") {
		return fmt.Errorf("invalid mongodb URI format")
	}

	// Validiere Port Ranges
	if config.API.Port < 1024 || config.API.Port > 65535 {
		return fmt.Errorf("api port must be between 1024 and 65535")
	}

	// Validiere Scraper Konfigurationen
	for name, sc := range config.Scrapers {
		if sc.DefaultPages <= 0 || sc.MaxPages <= 0 {
			return fmt.Errorf("invalid page configuration for scraper %s", name)
		}
		if sc.DefaultPages > sc.MaxPages {
			return fmt.Errorf("default_pages cannot be greater than max_pages for scraper %s", name)
		}

		// Validiere Schedule Format
		if _, err := cron.ParseStandard(sc.Schedule); err != nil {
			return fmt.Errorf("invalid schedule format for scraper %s: %w", name, err)
		}
	}

	return nil
}

// Hilfsfunktionen für Typkonvertierung
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func toInt(v interface{}, defaultValue int) int {
	if v == nil {
		return defaultValue
	}
	switch value := v.(type) {
	case int:
		return value
	case float64:
		return int(value)
	case string:
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getConfigValue(path string, rawValue interface{}) string {
	if rawValue == nil {
		return ""
	}

	configValue := fmt.Sprintf("%v", rawValue)
	if strings.HasPrefix(configValue, "${") && strings.HasSuffix(configValue, "}") {
		// Es ist ein Platzhalter, also Environment-Variable verwenden
		envVar := strings.TrimSuffix(strings.TrimPrefix(configValue, "${"), "}")
		if value := os.Getenv(envVar); value != "" {
			return value
		}
	}

	// Kein Platzhalter oder keine ENV var gefunden, config.yaml-Wert verwenden
	return viper.GetString(path)
}

// RequiredConfigError repräsentiert einen fehlenden Required Config Wert
type RequiredConfigError struct {
	Field string
}

func (e *RequiredConfigError) Error() string {
	return fmt.Sprintf("required configuration field missing: %s", e.Field)
}
