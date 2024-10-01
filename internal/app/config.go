package app

import (
	"job-scraper/internal/config"
)

func loadConfig() (*config.Config, error) {
	return config.LoadConfig()
}
