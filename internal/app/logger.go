package app

import (
	"job-scraper/internal/config"
	"job-scraper/internal/logging"
)

func initLogger(cfg *config.Config) {
	logging.InitLogger(cfg.Logging.Level)
}
