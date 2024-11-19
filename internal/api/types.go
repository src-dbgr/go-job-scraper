package api

// ScraperStatus repräsentiert den Status eines laufenden Scrapers
type ScraperStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Jobs   int    `json:"jobs"`
}
