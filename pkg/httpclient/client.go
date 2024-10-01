package httpclient

import (
	"net/http"
	"time"
)

func NewClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 30,
	}
}
