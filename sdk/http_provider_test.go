package sdk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPProvider_Load_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		_, _ = w.Write([]byte(`{"x":{"variants": {"yes":true, "no":false}, "defaultVariant":"yes"}}`))
	}))
	defer server.Close()

	p := NewHTTPProvider(server.URL, "", time.Second).(*httpProvider)
	store, err := p.Load(context.Background())
	assert.NoError(t, err)
	f, _ := store.Get("x")
	result := f.Evaluate(nil)
	assert.True(t, result.Value.(bool))
}

func TestHTTPProvider_Load_304(t *testing.T) {
	var called int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&called, 1) > 1 {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		_, _ = w.Write([]byte(`{"flags":{"y":{"enabled":false}}}`))
	}))
	defer server.Close()

	p := NewHTTPProvider(server.URL, "", time.Second).(*httpProvider)

	_, _ = p.Load(context.Background())        // first load: 200
	store, err := p.Load(context.Background()) // second load: 304

	assert.NoError(t, err)
	assert.Nil(t, store)
}

func TestHTTPProvider_Load_Error(t *testing.T) {
	p := NewHTTPProvider("https://invalid\\xZZ", "", time.Second).(*httpProvider)
	_, err := p.Load(context.Background())
	assert.Error(t, err)
}

func TestHTTPProvider_Load_400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	p := NewHTTPProvider(server.URL, "", time.Second).(*httpProvider)
	_, err := p.Load(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "http error")
}

func TestHTTPProvider_Load_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	p := NewHTTPProvider(server.URL, "", time.Second).(*httpProvider)
	_, err := p.Load(context.Background())
	assert.Error(t, err)
}

func TestHTTPProvider_Watch_OnlyFiresOnChange(t *testing.T) {
	var body atomic.Value
	body.Store(`{"feature":{"variants": {"yes":true, "no":false}, "defaultVariant":"yes"}}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		_, _ = w.Write([]byte(body.Load().(string)))
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	p := NewHTTPProvider(server.URL, "", 500*time.Millisecond).(*httpProvider)
	var hits int32

	// force first load
	_, _ = p.Load(context.Background())

	go p.Watch(ctx, func(s *Store) {
		if s != nil {
			f, _ := s.Get("feature")
			result := f.Evaluate(nil)
			if result.Value.(bool) {
				atomic.AddInt32(&hits, 1)
			}
		}
	})

	time.Sleep(1100 * time.Millisecond)
	assert.GreaterOrEqual(t, hits, int32(1))
}
