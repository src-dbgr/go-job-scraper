package parser

import (
	"encoding/json"
	"fmt"
	"job-scraper/internal/apperrors"
	"job-scraper/internal/models"
	"strconv"
	"time"
)

type JobParser struct{}

func NewJobParser() *JobParser {
	return &JobParser{}
}

func (jp *JobParser) ParseJob(data []byte) (*models.Job, error) {
	var aux struct {
		models.Job
		PostingDate       string      `json:"postingDate"`
		ExpirationDate    string      `json:"expirationDate"`
		YearsOfExperience interface{} `json:"yearsOfExperience"`
		CompanySize       interface{} `json:"companySize"`
		Remote            interface{} `json:"remote"`
		Salary            interface{} `json:"salary"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, apperrors.NewBaseError(apperrors.ErrCodeParser, "Failed to parse job data", err)
	}

	job := &aux.Job

	var err error
	job.PostingDate, err = parseDate(aux.PostingDate)
	if err != nil {
		return nil, apperrors.NewBaseError(apperrors.ErrCodeParser, "Error parsing posting date", err)
	}

	job.ExpirationDate, err = parseDate(aux.ExpirationDate)
	if err != nil {
		return nil, apperrors.NewBaseError(apperrors.ErrCodeParser, "Error parsing expiration date", err)
	}

	job.YearsOfExperience = parseIntOrString(aux.YearsOfExperience)
	job.CompanySize = parseIntOrString(aux.CompanySize)
	job.Remote = parseBool(aux.Remote)
	job.Salary = parseSalary(aux.Salary)

	return job, nil
}

func parseDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02",
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, apperrors.NewBaseError(
		apperrors.ErrCodeParser,
		fmt.Sprintf("Unable to parse date: %s", dateStr),
		nil,
	)
}

func parseIntOrString(v interface{}) int {
	switch value := v.(type) {
	case float64:
		return int(value)
	case string:
		if value == "" || value == "Not specified" {
			return 0
		}
		i, _ := strconv.Atoi(value)
		return i
	default:
		return 0
	}
}

func parseBool(v interface{}) bool {
	switch value := v.(type) {
	case bool:
		return value
	case string:
		return value == "Yes" || value == "true"
	default:
		return false
	}
}

func parseSalary(v interface{}) string {
	switch value := v.(type) {
	case string:
		return value
	case float64:
		return fmt.Sprintf("%.2f", value)
	case int:
		return strconv.Itoa(value)
	default:
		return ""
	}
}
