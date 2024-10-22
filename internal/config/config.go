package config

import (
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	API struct {
		Port int
	}
	MongoDB struct {
		URI      string
		Database string
	}
	Scrapers map[string]map[string]string
	OpenAI   struct {
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
	// Check for config path in environment
	if configPath := os.Getenv("JOBSCRAPER_CONFIG_PATH"); configPath != "" {
		viper.AddConfigPath(configPath)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("./configs") // default path
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("JOBSCRAPER")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}

	// API configuration
	apiPort, err := strconv.Atoi(getEnv("API_PORT", viper.GetString("api.port")))
	if err != nil {
		config.API.Port = 8080 // Default port
	} else {
		config.API.Port = apiPort
	}

	// MongoDB configuration
	config.MongoDB.URI = getEnv("MONGODB_URI", viper.GetString("mongodb.uri"))
	config.MongoDB.Database = getEnv("MONGODB_DATABASE", viper.GetString("mongodb.database"))

	// OpenAI configuration
	config.OpenAI.APIKey = getEnv("OPENAI_API_KEY", viper.GetString("openai.api_key"))
	config.OpenAI.APIURL = getEnv("OPENAI_API_URL", viper.GetString("openai.api_url"))
	config.OpenAI.Model = getEnv("OPENAI_MODEL", viper.GetString("openai.model"))
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
	prometheusPort, err := strconv.Atoi(getEnv("PROMETHEUS_PORT", viper.GetString("prometheus.port")))
	if err == nil {
		config.Prometheus.Port = prometheusPort
	} else {
		config.Prometheus.Port = viper.GetInt("prometheus.port")
	}

	// Scrapers configuration
	config.Scrapers = make(map[string]map[string]string)
	scrapers := viper.GetStringMap("scrapers")
	for scraperName, scraperConfig := range scrapers {
		config.Scrapers[scraperName] = make(map[string]string)
		if scraperConfigMap, ok := scraperConfig.(map[string]interface{}); ok {
			for key, value := range scraperConfigMap {
				config.Scrapers[scraperName][key] = getEnv("SCRAPER_"+scraperName+"_"+key, value.(string))
			}
		}
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
