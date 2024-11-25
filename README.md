# Job Scraper - Go Application

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8.svg)](https://golang.org/doc/devel/release.html)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

A scalable job scraping application written in Go that collects, processes, and analyzes job postings from various sources. 
The application uses by default the OpenAI ChatGPT API (or any other LLM Service that has API support) for intelligent data extraction and provides comprehensive analytics through Prometheus and Grafana.

<details>
<summary><strong>Table of Contents</strong></summary>

- [Job Scraper - Go Application](#job-scraper---go-application)
  - [Features](#features)
    - [Core Features](#core-features)
  - [Architecture](#architecture)
  - [Initial Setup and Configuration](#initial-setup-and-configuration)
    - [Prerequisites](#prerequisites)
    - [Environment Setup](#environment-setup)
    - [Deployment and Local exectuion Options](#deployment-and-local-exectuion-options)
      - [Local Development](#local-development)
      - [Starting/Running Locally](#startingrunning-locally)
      - [Docker Deployment](#docker-deployment)
      - [Kubernetes Deployment](#kubernetes-deployment)
      - [Service Access](#service-access)
    - [Verification](#verification)
  - [Usage](#usage)
    - [API Endpoints](#api-endpoints)
      - [Scraping Operations](#scraping-operations)
      - [Data Access](#data-access)
  - [Monitoring \& Observability](#monitoring--observability)
    - [Prometheus Metrics](#prometheus-metrics)
      - [API Metrics](#api-metrics)
      - [Scraper Metrics](#scraper-metrics)
      - [Processor Metrics](#processor-metrics)
      - [Storage Metrics](#storage-metrics)
    - [Prometheus Configuration](#prometheus-configuration)
    - [Grafana Setup and Usage](#grafana-setup-and-usage)
      - [Initial Access](#initial-access)
      - [Setting up the Grafana API Key to Set Up the Data Sources](#setting-up-the-grafana-api-key-to-set-up-the-data-sources)
      - [Installing the JSON API Connection to Grafana](#installing-the-json-api-connection-to-grafana)
      - [Automated Datasource Setup using Scripts](#automated-datasource-setup-using-scripts)
      - [Importing the Dashboard](#importing-the-dashboard)
      - [Additional Data Sources (Optional)](#additional-data-sources-optional)
      - [Dashboard Maintenance](#dashboard-maintenance)
    - [Troubleshooting](#troubleshooting)
  - [Extending the Application](#extending-the-application)
    - [Adding a New Scraper](#adding-a-new-scraper)
    - [Adding New Metrics](#adding-new-metrics)
    - [Adding New API Endpoints](#adding-new-api-endpoints)
  - [Testing](#testing)
    - [Running Tests](#running-tests)
    - [Code Style Guidelines](#code-style-guidelines)
    - [Security Considerations](#security-considerations)
  - [License](#license)
</details>

## Features

### Core Features
- Modular scraper architecture extensible for multiple job portals
- Intelligent job data extraction using ChatGPT
- MongoDB persistence layer
- RESTful API for data access and control
- Comprehensive metrics and monitoring
- Kubernetes-ready deployment
- Scheduled scraping with configurable intervals

## Architecture

```mermaid
%%{init: {'theme': 'base', 'themeVariables': { 'lineColor': '#000' }}}%%
graph TD
    subgraph Architecture ["System Architecture"]
        subgraph External [" "]
            A[REST API]
            CR[⏰ Cron Scheduler]
        end
        
        subgraph Core ["Core Application"]
            A --> B[Scraper Service]
            CR --> B
            B --> C[Job Processors]
            C --> D[Storage Layer]
            B --> E[Metrics Collector]
            C --> F[LLM Processor Interface]
            D --> G[(MongoDB)]
            E --> H[Prometheus]
            
            subgraph LLM ["LLM Services"]
                F --> I[OpenAI Service]
                F --> J[...]
                F --> K[Other LLM API Services, i.e. Anthropic, Gemini etc.]
            end
        end

        GF([📊 Grafana Dashboard]) -.-> H
        GF -.-> G
    end
    
    style A fill:#fff,stroke:#000,color:#000
    style CR fill:#fff,stroke:#000,color:#000
    style B fill:#fff,stroke:#000,color:#000
    style C fill:#fff,stroke:#000,color:#000
    style D fill:#fff,stroke:#000,color:#000
    style E fill:#fff,stroke:#000,color:#000
    style F fill:#fff,stroke:#000,color:#000
    style G fill:#fff,stroke:#000,color:#000
    style H fill:#fff,stroke:#000,color:#000
    style I fill:#fff,stroke:#000,color:#000
    style GF fill:#fff,stroke:#000,color:#000

    style J fill:#f0f0f0,stroke:#000,color:#000
    style K fill:#f0f0f0,stroke:#000,color:#000

    classDef subgraphStyle fill:transparent,stroke-dasharray: 5 5,stroke:#000,color:#000;
    class LLM subgraphStyle;

    classDef coreStyle fill:#fafafa,stroke-dasharray: 2 2, stroke:#000,color:#000;
    class Core,External coreStyle;

    classDef architectureStyle fill:#fff,stroke:#000,color:#000;
    class Architecture architectureStyle;
```

The application follows a modular, layered architecture:
- **API Layer**: HTTP endpoints for control and data access
- **Service Layer**: Business logic and orchestration
- **Processor Layer**: Data processing and enrichment
- **Storage Layer**: Data persistence and retrieval
- **Metrics Layer**: Performance and operational metrics

## Initial Setup and Configuration

### Prerequisites
- Go 1.23+
- Docker and Docker Compose
- MongoDB 6.0+
- OpenAI API key
- Kubernetes (for production deployment)

### Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/src-dbgr/go-job-scraper.git
cd go-job-scraper
```

2. Configure environment variables:
```bash
# Create an environment file
touch .env
```

3. Set required environment variables in `.env`:

Variables Explanation:
- MONGODB_URI: MongoDB connection string
- MONGODB_DATABASE: Database name
- OPENAI_API_KEY: Your OpenAI API key
- SCRAPER_JOBSCH_BASE_URL: Base URL for the jobs.ch API
- SCRAPER_JOBSCH_API_KEY: API key for jobs.ch (if required)
- GRAFANA_ADMIN_PASSWORD: Password for Grafana admin user

Example
```bash
# Core Settings
MONGODB_URI=mongodb://mongodb:27017
MONGODB_DATABASE=jobsdb
OPENAI_API_KEY=some-dummy-key
OPENAI_API_URL=https://api.openai.com/v1/chat/completions

# Monitoring Settings
GRAFANA_ADMIN_PASSWORD=admin
```

4. Configure application settings if required in `config.yaml`:
```yaml
api:
  port: 8080  # Default port, overwritten by env var

mongodb:
  uri: ${MONGODB_URI}
  database: ${MONGODB_DATABASE}

scrapers:
  jobsch:
    base_url: https://www.jobs.ch/api/v1
    api_key: ${SCRAPER_JOBSCH_API_KEY}
    default_pages: 5         # Default # of pages for the scheduler
    max_pages: 20            # No. of jobs per page
    schedule: "0 */6 * * *"  # Cron expression for every 6 hours

logging:
  level: "info"
  file: "logs/job_scraper.log"

prometheus:
  port: 2112

openai:
  api_key: ${OPENAI_API_KEY}
  api_url: ${OPENAI_API_URL}
  model: gpt-4o-mini
  timeout: 300s
  temperature: 1
  max_tokens: 500
  top_p: 1
  frequency_penalty: 0
  presence_penalty: 0
```

### Deployment and Local exectuion Options

#### Local Development
1. Install dependencies:
```bash
make deps
```

2. Start required services:
```bash
docker compose up -d mongodb prometheus grafana
```

3. Build and run the application:
```bash
make build
./dist/job-scraper
```

#### Starting/Running Locally
1. cd into `go-job-scraper`
2. Start the application by executing `go run ./cmd/scraper/main.go`
3. Trigger the scrape process: `curl -X POST "http://localhost:8080/api/v1/scrape/jobsch?pages=2"`
   1. Adjust the number of considered pages to your needs. Notice, you need to have a ChatGPT API Key in place.
   2. One page on jobsch contains 20 Jobs, so 2 pages as set here will contain 40 jobs that will be processed by ChatGPT

#### Docker Deployment
Start the complete application stack:
```bash
docker compose up -d
```

This includes:
- Job Scraper application
- MongoDB database
- Prometheus monitoring
- Grafana dashboards

#### Kubernetes Deployment
1. Create the namespace:
```bash
kubectl create namespace job-scraper
```

2. Create the secrets (replace values with your encoded secrets):
```bash
kubectl create -f deployments/kubernetes/secrets.yaml
```

3. Deploy using kustomize:
```bash
kubectl apply -k deployments/kubernetes/
```

4. Verify the deployment:
```bash
kubectl -n job-scraper get pods
kubectl -n job-scraper get services
kubectl -n job-scraper get ingress
```

Note: Before deploying, make sure to:
- Replace `your-registry` with your actual container registry
- Update the host in the Ingress configuration
- Set your actual secrets in `secrets.yaml`
- Adjust resource limits based on your needs

#### Service Access
After successful deployment, access services at:
- Job Scraper API: http://localhost:8080
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000
- MongoDB: localhost:27017

### Verification
Verify the setup with:
```bash
# Check service health
curl http://localhost:8080/health

# Verify metrics endpoint
curl http://localhost:2112/metrics

# Test API access
curl http://localhost:8080/api/v1/scrapers/status
```

## Usage

### API Endpoints

#### Scraping Operations
```bash
# Start a scraping job
curl -X POST http://localhost:8080/api/v1/scrape/jobsch

# Check scraper status
curl http://localhost:8080/api/v1/scrapers/status

# Get specific job by ID
curl http://localhost:8080/api/v1/jobs/{id}
```

#### Data Access
```bash
# Get all jobs
curl http://localhost:8080/api/v1/jobs

# Get job statistics
curl http://localhost:8080/api/v1/stats/job-categories-counts
```

## Monitoring & Observability

### Prometheus Metrics

The application exposes metrics at `:2112/metrics`. Here are the key metric categories:

#### API Metrics
```go
HTTPRequestDuration // Duration of HTTP requests
HTTPRequestsTotal   // Total number of HTTP requests
ActiveRequests      // Number of currently active requests
```

#### Scraper Metrics
```go
ScrapingDuration    // Duration of scraping operations
ScrapedJobsTotal    // Total number of scraped jobs
ScraperErrors       // Total number of scraper errors
```

#### Processor Metrics
```go
ProcessorDuration   // Duration of processing operations
ProcessorErrors     // Total number of processor errors
OpenAITokensUsed    // Total number of OpenAI tokens used
```

#### Storage Metrics
```go
DBOperationDuration // Duration of database operations
DBOperationsTotal   // Total number of database operations
DBConnectionsActive // Number of active database connections
```

### Prometheus Configuration

The application uses the following Prometheus config (prometheus.yml):
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'job-scraper'
    static_configs:
      - targets: ['host.docker.internal:2112']
```

### Grafana Setup and Usage

#### Initial Access
1. Access Grafana at http://localhost:3000
2. Login with default credentials:
   - Username: `admin`
   - Password: Value from `GRAFANA_ADMIN_PASSWORD` in your `.env`

#### Setting up the Grafana API Key to Set Up the Data Sources
Before using the automation scripts, you need to create a Grafana API key:

1. Go to Administration > Users and access > Service accounts
2. Click "Add service account"
3. Display name *: `job-scraper-admin`
4. Role: `Admin`
5. Click "Create"
6. Click "Add service account token"
7. Add Display name *: `job-scraper-admin`
8. Choose the Expiration as you like
9. Click "Generate token"
10. Copy the token by clicking: "Copy to clipboard and close"
11. **Important**: Save the generated API key securely, you will need it for the datasource setup script

#### Installing the JSON API Connection to Grafana
1. Go to Connections > Add new connection
2. Search for `json api`
3. Choose `JSON API` result
4. Install

#### Automated Datasource Setup using Scripts

The project includes an automation script to set up all required datasources:

```bash
# Make the script executable
chmod +x scripts/manage_grafana_v1_datasources.sh

# Run the script (you will be prompted for your Grafana API key)
./scripts/manage_grafana_v1_datasources.sh
# Choose option 1 when prompted to create all datasources
```

This script will create all necessary datasources including:
- Jobs endpoint
- Job category statistics
- Salary statistics
- Various job market analytics endpoints

#### Importing the Dashboard

The project includes a pre-configured dashboard for job scraping analytics:

1. Go to Dashboards > New > Import
2. Click "Upload JSON file"
3. Select `configs/dashboards/job_scaper_dashboard.json`
4. Click "Import"

The dashboard includes:
- Job scraping overview
- Category distribution
- Geographical insights
- Salary trends
- Processing metrics

> Note: Your scraping application must be running as it exposes the necessary REST endpoints. Also you need scraped data to be available in your MongoDB in order to visualize anything in grafana

#### Additional Data Sources (Optional)

While the main dashboard uses the JSON API datasources, you can optionally add Prometheus as an additional data source for system metrics:

1. Go to Configuration > Data Sources
2. Add Prometheus
3. URL: `http://prometheus:9090`
4. Click "Save & Test"

This allows monitoring of system-level metrics like scraping performance and API response times alongside the job market analytics.

#### Dashboard Maintenance

Check the dashboard's health:
1. Verify all panels are loading data correctly
2. Check API endpoint connectivity
3. Review any error messages in panel queries
4. Update time ranges for relevant insights

If you need to reset or recreate the datasources, you can use the script with the delete option:
```bash
./scripts/manage_grafana_v1_datasources.sh
# Choose option 2 when prompted to delete all datasources
```

### Troubleshooting

Common issues and solutions:

1. Metrics not showing up:
   - Check if the application is exposing metrics on port 2112 for Prometheus or 8080 for API access, example: 
      ```bash
      curl http://localhost:8080/api/v1/stats/job-categories-counts
      ```
   - Verify Prometheus target is reachable
   - Check for any firewall issues

2. Grafana connection issues:
   - Verify Prometheus data source configuration (according Prometheus datasource needs to be configured) and the Datasource connection URL needs to be: `http://host.docker.internal:<port>` in case you run it locally in Docker
   - Check network connectivity between containers
   - Validate authentication settings
   - MongoDB must be up and running and should contain Documents (scraped job data) in database `jobsdb` that hold the collection `jobs`
     - You can utilize a tool such as `MongoDB Compass` to inspect the database
     
     - Checking MongoDB connection with mongosh:
      ```bash
      # Using mongosh
      mongosh "mongodb://localhost:27017/jobsdb"

      # Check collections
      show collections

      # Check job data
      db.jobs.findOne()
      ```

3. Missing data points:
   - Check scrape interval configuration
   - Verify metric collection is active
   - Check for any rate limiting issues

## Extending the Application

### Adding a New Scraper

1. Create a new scraper package:
```go
// internal/scraper/newportal/newportal_scraper.go

package newportal

type NewPortalScraper struct {
    client     HTTPClient
    baseURL    string
    jobFetcher JobFetcher
}

func NewNewPortalScraper(config Config) *NewPortalScraper {
    return &NewPortalScraper{
        client:  &http.Client{},
        baseURL: config.BaseURL,
    }
}

func (s *NewPortalScraper) Scrape(ctx context.Context) ([]models.Job, error) {
    // Implement scraping logic
}
```

2. Register the scraper in the factory:
```go
// internal/scraper/factory.go

func NewScraper(name string, config map[string]string) (Scraper, error) {
    switch name {
    case "newportal":
        return newNewPortalScraper(config)
    // ...
    }
}
```

3. Add configuration:
```yaml
scrapers:
  newportal:
    base_url: https://api.newportal.com
    schedule: "0 */6 * * *"
```

### Adding New Metrics

1. Define metrics in a domain file:
```go
// internal/metrics/domains/newmetrics.go

var (
    NewMetric = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "jobscraper_new_metric_total",
            Help: "Description of the metric",
        },
    )
)
```

2. Implement collection logic in relevant components.

### Adding New API Endpoints

1. Create handler function:
```go
// internal/api/new_handler.go

func (a *API) handleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Implement handler logic
}
```

2. Register route:
```go
// internal/api/api.go

func (a *API) setupRoutes() {
    v1Router.HandleFunc("/new/endpoint", a.handleNewEndpoint).Methods("GET")
}
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Generate coverage report
make cover
```

### Code Style Guidelines

- Follow Go best practices and idioms
- Use meaningful variable and function names
- Write tests for new functionality
- Update documentation as needed
- Add appropriate logging and metrics

### Security Considerations
- All API keys and sensitive data should be stored securely
- The application doesn't implement authentication by default
- For production deployments, consider adding:
  - API authentication
  - HTTPS/TLS
  - Network policies in Kubernetes
  - Rate limiting

## License

This project is licensed under the GPL v3 License - see the [LICENSE](LICENSE) file for details.
