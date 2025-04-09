package sdk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type URLStoreOption func(r *http.Request)

var NotModifiedError = errors.New("not modified")

// SetLastModified is an URLStoreOption for setting a last modified date, which can result in a 304
func SetLastModified(modTime time.Time) URLStoreOption {
	return func(r *http.Request) {
		r.Header.Set("If-Modified-Since", modTime.Format(http.TimeFormat))
	}
}

// NewStoreFromURL loads a flag file from an HTTP(S) endpoint.
// It auto-detects JSON or YAML based on the file extension.
func NewStoreFromURL(url string, token string, options ...URLStoreOption) (*Store, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Apply our extended options
	for _, opt := range options {
		opt(req)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch flags: %w", err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Not modified since last request
	if resp.StatusCode == http.StatusNotModified {
		return nil, NotModifiedError
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	// error will come out further down anyway (and is easier to test)
	body, _ := io.ReadAll(resp.Body)

	format := DetectFormat(url)
	return NewStoreFromBytesWithFormat(body, format)
}
