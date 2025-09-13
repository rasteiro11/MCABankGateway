package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type requestOptions struct {
	baseURL string
	timeout time.Duration
	headers map[string]string
	query   url.Values
}

func DoRequest[T any](ctx context.Context, method, path string, body any, opts ...Option) (*T, error) {
	ro := &requestOptions{
		timeout: 10 * time.Second,
	}
	for _, o := range opts {
		o(ro)
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	fullURL := ro.baseURL + path
	if len(ro.query) > 0 {
		fullURL = fullURL + "?" + ro.query.Encode()
	}

	client := &http.Client{Timeout: ro.timeout}
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range ro.headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return nil, NewHTTPError(resp.StatusCode, string(data))
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func Get[T any](ctx context.Context, path string, opts ...Option) (*T, error) {
	return DoRequest[T](ctx, http.MethodGet, path, nil, opts...)
}

func Post[T any](ctx context.Context, path string, body any, opts ...Option) (*T, error) {
	return DoRequest[T](ctx, http.MethodPost, path, body, opts...)
}

func Put[T any](ctx context.Context, path string, body any, opts ...Option) (*T, error) {
	return DoRequest[T](ctx, http.MethodPut, path, body, opts...)
}

func Delete[T any](ctx context.Context, path string, opts ...Option) (*T, error) {
	return DoRequest[T](ctx, http.MethodDelete, path, nil, opts...)
}
