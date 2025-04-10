package sdk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// httpProvider implements StoreProvider by polling an HTTP endpoint.
type httpProvider struct {
	URL       string
	Token     string
	Interval  time.Duration
	lastMod   string
	lastStore *Store
	mu        sync.Mutex
}

func NewHTTPProvider(url string, token string, interval time.Duration) StoreProvider {
	return &httpProvider{
		URL:      url,
		Token:    token,
		Interval: interval,
	}
}

func (p *httpProvider) Load(ctx context.Context) (*Store, error) {

	// Construct the request
	req, err := http.NewRequestWithContext(ctx, "GET", p.URL, nil)
	if err != nil {
		return nil, err
	}
	if p.Token != "" {
		req.Header.Set("Authorization", "Bearer "+p.Token)
	}
	if p.lastMod != "" {
		req.Header.Set("If-Modified-Since", p.lastMod)
	}

	// Do the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// 304 = no change
	if resp.StatusCode == http.StatusNotModified {
		return nil, nil
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	store, err := NewStoreFromBytesWithFormat(body, DetectFormat(p.URL))
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastStore = store
	p.lastMod = resp.Header.Get("Last-Modified")
	return store, nil
}

func (p *httpProvider) Watch(ctx context.Context, onChange func(*Store)) {
	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			store, err := p.Load(ctx)
			if err != nil || store == nil {
				continue
			}
			onChange(store)
		}
	}
}
