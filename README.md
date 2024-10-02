# Job Scraper

A job scraping application written in Go

## Features

- Scrapes job listings from various sources
- Processes job descriptions using ChatGPT
- Stores data in MongoDB
- Provides metrics via Prometheus
- Deployable on Kubernetes

## Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Kubernetes (Optional)

## Local Development

1. Clone the repository:
   ```
   git clone https://github.com/src-dbgr/go-job-scraper.git
   cd go-job-scraper
   ```

2. Set up environment variables:
   ```
   export OPENAI_API_KEY=your_chatgpt_api_key
   export OPENAI_API_URL=your_chatgpt_api_url
   export GRAFANA_ADMIN_PASSWORD=your_grafana_password
   ```

3. Start the application using Docker Compose:
   ```
   docker-compose up --build
   ```

4. Access the services:
   - Job Scraper metrics: http://localhost:2112/metrics
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3000

## Running Tests

```
make test
```

## Linting

```
make lint
```

## Deployment

1. Build and push the Docker image:
   ```
   docker build -t your-registry/job-scraper:latest .
   docker push your-registry/job-scraper:latest
   ```

2. Update Kubernetes manifests in `deployments/kubernetes/` with your image and secrets.

3. Apply the Kubernetes manifests:
   ```
   kubectl apply -f deployments/kubernetes/
   ```

## ... To Be Continued