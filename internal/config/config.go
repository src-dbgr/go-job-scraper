package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDB struct {
		URI      string
		Database string
	}
	Scrapers map[string]map[string]string
	ChatGPT  struct {
		APIKey string
		APIURL string
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
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("JOBSCRAPER")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}

	// MongoDB configuration
	config.MongoDB.URI = getEnv("MONGODB_URI", viper.GetString("mongodb.uri"))
	config.MongoDB.Database = getEnv("MONGODB_DATABASE", viper.GetString("mongodb.database"))

	// ChatGPT configuration
	config.ChatGPT.APIKey = getEnv("CHATGPT_API_KEY", viper.GetString("chatgpt.api_key"))
	config.ChatGPT.APIURL = getEnv("CHATGPT_API_URL", viper.GetString("chatgpt.api_url"))

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
