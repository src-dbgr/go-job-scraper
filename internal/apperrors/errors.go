package apperrors

import (
	"fmt"
)

const (
	// Error Codes für verschiedene Kategorien
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeScraping       = "SCRAPING_ERROR"
	ErrCodeProcessing     = "PROCESSING_ERROR"
	ErrCodeStorage        = "STORAGE_ERROR"
	ErrCodeConfig         = "CONFIG_ERROR"
	ErrCodeInitialization = "INITIALIZATION_ERROR"
	ErrCodeMetrics        = "METRICS_ERROR"
	ErrCodeScheduler      = "SCHEDULER_ERROR"
	ErrCodeParser         = "PARSER_ERROR"
	ErrCodeInternal       = "INTERNAL_ERROR"
)

// BaseError ist der Basis-Fehlertyp
type BaseError struct {
	Code    string
	Message string
	Err     error
}

func NewBaseError(code string, message string, err error) *BaseError {
	return &BaseError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *BaseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *BaseError) Unwrap() error {
	return e.Err
}

// NotFoundError für nicht gefundene Ressourcen
type NotFoundError struct {
	*BaseError
	Resource string
	ID       string
}

func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		BaseError: NewBaseError(ErrCodeNotFound, fmt.Sprintf("%s with ID %s not found", resource, id), nil),
		Resource:  resource,
		ID:        id,
	}
}

// ScrapingError für Scraping-Fehler
type ScrapingError struct {
	*BaseError
	Source string
	URL    string
}

func NewScrapingError(source, url string, err error) *ScrapingError {
	return &ScrapingError{
		BaseError: NewBaseError(ErrCodeScraping, fmt.Sprintf("failed to scrape from %s", source), err),
		Source:    source,
		URL:       url,
	}
}

// ProcessingError für Verarbeitungsfehler
type ProcessingError struct {
	*BaseError
	JobID string
}

func NewProcessingError(jobID string, message string, err error) *ProcessingError {
	return &ProcessingError{
		BaseError: NewBaseError(ErrCodeProcessing, message, err),
		JobID:     jobID,
	}
}
