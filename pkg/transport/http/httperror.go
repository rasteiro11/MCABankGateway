package http

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error: %d %s - %s", e.StatusCode, e.Status, e.Body)
}

func NewHTTPError(statusCode int, body string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       body,
	}
}
