package app

import (
	"context"
	"job-scraper/internal/config"
	"job-scraper/internal/scheduler"
	"job-scraper/internal/scraper"
	"job-scraper/internal/services"

	"github.com/robfig/cron/v3"
)

func initScheduler(
	ctx context.Context,
	scraperService *services.ScraperService,
	scrapers map[string]scraper.Scraper,
	cfg *config.Config,
) (*scheduler.Scheduler, error) {
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	var scraperConfigs []scheduler.ScraperConfig
	for scraperType, scraperCfg := range cfg.Scrapers {
		schedule, err := cronParser.Parse(scraperCfg.Schedule)
		if err != nil {
			return nil, err
		}
		scraperConfigs = append(scraperConfigs, scheduler.ScraperConfig{
			Type:     scraperType,
			Schedule: schedule,
			Pages:    scraperCfg.DefaultPages,
		})
	}

	return scheduler.NewScheduler(scraperService, scrapers, scraperConfigs), nil
}
