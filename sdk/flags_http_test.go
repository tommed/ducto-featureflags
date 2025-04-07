package sdk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStoreFromURL_JSON(t *testing.T) {
	payload := `{
		"flags": {
			"new_ui": {
				"enabled": true
			}
		}
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()

	store, err := NewStoreFromURL(srv.URL+"/flags", "test-token")
	assert.NoError(t, err)
	assert.True(t, store.IsEnabled("new_ui", EvalContext{}))
}

func TestNewStoreFromURL_YAML(t *testing.T) {
	payload := `
flags:
  canary:
    enabled: false
    rules:
      - if:
          env: prod
        value: true
`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()

	store, err := NewStoreFromURL(srv.URL+"/flags.yaml", "")
	assert.NoError(t, err)

	assert.True(t, store.IsEnabled("canary", EvalContext{"env": "prod"}))
	assert.False(t, store.IsEnabled("canary", EvalContext{"env": "dev"}))
}
