package http

import (
	"net/url"
	"time"
)

type Option func(*requestOptions)

func WithBaseURL(url string) Option {
	return func(ro *requestOptions) {
		ro.baseURL = url
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(ro *requestOptions) {
		ro.timeout = timeout
	}
}

func WithHeader(key, value string) Option {
	return func(ro *requestOptions) {
		if ro.headers == nil {
			ro.headers = make(map[string]string)
		}
		ro.headers[key] = value
	}
}

func WithQueryParam(key, value string) Option {
	return func(ro *requestOptions) {
		if ro.query == nil {
			ro.query = make(url.Values)
		}
		ro.query.Add(key, value)
	}
}

func WithQueryParams(params map[string]string) Option {
	return func(ro *requestOptions) {
		if ro.query == nil {
			ro.query = make(url.Values)
		}
		for k, v := range params {
			ro.query.Add(k, v)
		}
	}
}
