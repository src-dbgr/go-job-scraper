package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"job-scraper/internal/api/middleware"
	"job-scraper/internal/processor"
	"job-scraper/internal/scraper"
	"job-scraper/internal/services"
	"job-scraper/internal/storage"

	"github.com/gorilla/mux"
)

type API struct {
	router          *mux.Router
	scrapers        map[string]scraper.Scraper
	storage         storage.Storage
	processor       processor.JobProcessor
	scraperService  *services.ScraperService
	runningScrapers *sync.Map
	jobStatsService *services.JobStatisticsService
}

func NewAPI(
	scrapers map[string]scraper.Scraper,
	storage storage.Storage,
	processor processor.JobProcessor,
	scraperService *services.ScraperService,
	jobStatsService *services.JobStatisticsService,
) *API {
	api := &API{
		router:          mux.NewRouter(),
		scrapers:        scrapers,
		storage:         storage,
		processor:       processor,
		scraperService:  scraperService,
		runningScrapers: &sync.Map{},
		jobStatsService: jobStatsService,
	}
	api.setupRoutes()
	return api
}

func (a *API) setupRoutes() {

	a.router.Use(middleware.ErrorHandler)

	// Create a subrouter for v1
	v1Router := a.router.PathPrefix("/api/v1").Subrouter()

	// Middleware that sets the api version and metrics
	v1Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			VersionMiddleware(next.ServeHTTP, "v1")(w, r)
		})
	})
	v1Router.Use(middleware.MetricsMiddleware)

	// Scraper routes
	v1Router.HandleFunc("/scrape/{scraper}", a.handleScrape).Methods("POST")
	v1Router.HandleFunc("/scrapers/status", a.handleScrapersStatus).Methods("GET")

	// Job routes
	v1Router.HandleFunc("/jobs", a.getJobs).Methods("GET")
	v1Router.HandleFunc("/jobs/{id}", a.getJobByID).Methods("GET")
	v1Router.HandleFunc("/jobs/urls", a.getJobUrls).Methods("GET")

	// Statistics routes
	v1Router.HandleFunc("/stats/top-job-categories", a.getTopJobCategories).Methods("GET")
	v1Router.HandleFunc("/stats/avg-experience-by-category", a.getAvgExperienceByCategory).Methods("GET")
	v1Router.HandleFunc("/stats/remote-vs-onsite", a.getRemoteVsOnsite).Methods("GET")
	v1Router.HandleFunc("/stats/top-skills", a.getTopSkills).Methods("GET")
	v1Router.HandleFunc("/stats/top-optional-skills", a.getTopOptionalSkills).Methods("GET")
	v1Router.HandleFunc("/stats/benefits-by-company-size", a.getBenefitsByCompanySize).Methods("GET")
	v1Router.HandleFunc("/stats/avg-salary-by-education", a.getAvgSalaryByEducation).Methods("GET")
	v1Router.HandleFunc("/stats/job-postings-trend", a.getJobPostingsTrend).Methods("GET")
	v1Router.HandleFunc("/stats/languages-by-location", a.getLanguagesByLocation).Methods("GET")
	v1Router.HandleFunc("/stats/employment-types", a.getEmploymentTypes).Methods("GET")
	v1Router.HandleFunc("/stats/remote-work-by-category", a.getRemoteWorkByCategory).Methods("GET")
	v1Router.HandleFunc("/stats/technology-trends", a.getTechnologyTrends).Methods("GET")
	v1Router.HandleFunc("/stats/job-requirements-by-location", a.getJobRequirementsByLocation).Methods("GET")
	v1Router.HandleFunc("/stats/remote-vs-onsite-by-industry", a.getRemoteVsOnsiteByIndustry).Methods("GET")
	v1Router.HandleFunc("/stats/job-categories-by-company-size", a.getJobCategoriesByCompanySize).Methods("GET")
	v1Router.HandleFunc("/stats/skills-by-experience-level", a.getSkillsByExperienceLevel).Methods("GET")
	v1Router.HandleFunc("/stats/companies-by-size", a.getCompaniesBySize).Methods("GET")
	v1Router.HandleFunc("/stats/companies-by-size/{sizeType}", a.getCompaniesBySizeAndType).Methods("GET")
	v1Router.HandleFunc("/stats/company-size-distribution", a.getCompanySizeDistribution).Methods("GET")
	v1Router.HandleFunc("/stats/job-postings-per-day", a.getJobPostingsPerDay).Methods("GET")
	v1Router.HandleFunc("/stats/job-postings-per-month", a.getJobPostingsPerMonth).Methods("GET")
	v1Router.HandleFunc("/stats/job-postings-per-company", a.getJobPostingsPerCompany).Methods("GET")
	v1Router.HandleFunc("/stats/mustskills/{skill}", a.getMustSkillFrequencyPerDay).Methods("GET")
	v1Router.HandleFunc("/stats/optionalskills/{skill}", a.getOptionalSkillFrequencyPerDay).Methods("GET")
	v1Router.HandleFunc("/stats/job-categories-counts", a.getJobCategoryCounts).Methods("GET")

	a.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
