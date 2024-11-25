# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o job-scraper ./cmd/scraper

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Set env vars for the prompt path
ENV JOBSCRAPER_PROMPT_PATH=/app/prompts
ENV JOBSCRAPER_CONFIG_PATH=/app/configs

WORKDIR /app

COPY --from=builder /app/job-scraper .
COPY --from=builder /app/configs/ ./configs/
COPY --from=builder /app/prompts/ ./prompts/

# Prometheus port
EXPOSE 2112

# Actual application port
EXPOSE 8080

CMD ["./job-scraper"]