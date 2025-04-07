package sdk

import (
	"fmt"
	"io"
	"net/http"
)

// NewStoreFromURL loads a flag file from an HTTP(S) endpoint.
// It auto-detects JSON or YAML based on the file extension.
func NewStoreFromURL(url string, token string) (*Store, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch flags: %w", err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	format := detectFormat(url)
	return NewStoreFromBytesWithFormat(body, format)
}
