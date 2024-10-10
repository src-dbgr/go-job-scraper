package api

import (
	"net/http"
	"sync"

	"job-scraper/internal/processor/openai"
	"job-scraper/internal/scraper"
	"job-scraper/internal/services"
	"job-scraper/internal/storage"

	"github.com/gorilla/mux"
)

type API struct {
	router          *mux.Router
	scrapers        map[string]scraper.Scraper
	storage         storage.Storage
	openaiProcessor *openai.Processor
	runningScrapers *sync.Map
	jobStatsService *services.JobStatisticsService
}

type ScraperStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Jobs   int    `json:"jobs"`
}

func NewAPI(scrapers map[string]scraper.Scraper, storage storage.Storage, openaiProcessor *openai.Processor, jobStatsService *services.JobStatisticsService) *API {
	api := &API{
		router:          mux.NewRouter(),
		scrapers:        scrapers,
		storage:         storage,
		openaiProcessor: openaiProcessor,
		runningScrapers: &sync.Map{},
		jobStatsService: jobStatsService,
	}
	api.setupRoutes()
	return api
}

func (a *API) setupRoutes() {
	// Scraper routes
	a.router.HandleFunc("/api/scrape/{scraper}", a.handleScrape).Methods("POST")
	a.router.HandleFunc("/api/scrapers/status", a.handleScrapersStatus).Methods("GET")

	// Job routes
	a.router.HandleFunc("/api/jobs", a.getJobs).Methods("GET")
	a.router.HandleFunc("/api/jobs/{id}", a.getJobByID).Methods("GET")
	a.router.HandleFunc("/api/jobs/urls", a.getJobUrls).Methods("GET")

	// Statistics routes
	a.router.HandleFunc("/api/stats/top-job-categories", a.getTopJobCategories).Methods("GET")
	a.router.HandleFunc("/api/stats/avg-experience-by-category", a.getAvgExperienceByCategory).Methods("GET")
	a.router.HandleFunc("/api/stats/remote-vs-onsite", a.getRemoteVsOnsite).Methods("GET")
	a.router.HandleFunc("/api/stats/top-skills", a.getTopSkills).Methods("GET")
	a.router.HandleFunc("/api/stats/top-optional-skills", a.getTopOptionalSkills).Methods("GET")
	a.router.HandleFunc("/api/stats/benefits-by-company-size", a.getBenefitsByCompanySize).Methods("GET")
	a.router.HandleFunc("/api/stats/avg-salary-by-education", a.getAvgSalaryByEducation).Methods("GET")
	a.router.HandleFunc("/api/stats/job-postings-trend", a.getJobPostingsTrend).Methods("GET")
	a.router.HandleFunc("/api/stats/languages-by-location", a.getLanguagesByLocation).Methods("GET")
	a.router.HandleFunc("/api/stats/employment-types", a.getEmploymentTypes).Methods("GET")
	a.router.HandleFunc("/api/stats/remote-work-by-category", a.getRemoteWorkByCategory).Methods("GET")
	a.router.HandleFunc("/api/stats/technology-trends", a.getTechnologyTrends).Methods("GET")
	a.router.HandleFunc("/api/stats/job-requirements-by-location", a.getJobRequirementsByLocation).Methods("GET")
	a.router.HandleFunc("/api/stats/remote-vs-onsite-by-industry", a.getRemoteVsOnsiteByIndustry).Methods("GET")
	a.router.HandleFunc("/api/stats/job-categories-by-company-size", a.getJobCategoriesByCompanySize).Methods("GET")
	a.router.HandleFunc("/api/stats/skills-by-experience-level", a.getSkillsByExperienceLevel).Methods("GET")
	a.router.HandleFunc("/api/stats/companies-by-size", a.getCompaniesBySize).Methods("GET")
	a.router.HandleFunc("/api/stats/companies-by-size/{sizeType}", a.getCompaniesBySizeAndType).Methods("GET")
	a.router.HandleFunc("/api/stats/company-size-distribution", a.getCompanySizeDistribution).Methods("GET")
	a.router.HandleFunc("/api/stats/job-postings-per-day", a.getJobPostingsPerDay).Methods("GET")
	a.router.HandleFunc("/api/stats/job-postings-per-month", a.getJobPostingsPerMonth).Methods("GET")
	a.router.HandleFunc("/api/stats/job-postings-per-company", a.getJobPostingsPerCompany).Methods("GET")
	a.router.HandleFunc("/api/stats/mustskills/{skill}", a.getMustSkillFrequencyPerDay).Methods("GET")
	a.router.HandleFunc("/api/stats/optionalskills/{skill}", a.getOptionalSkillFrequencyPerDay).Methods("GET")
	a.router.HandleFunc("/api/stats/job-categories-counts", a.getJobCategoryCounts).Methods("GET")
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
