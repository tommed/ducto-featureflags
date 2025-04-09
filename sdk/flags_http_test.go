package sdk

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStoreFromURL_JSON(t *testing.T) {
	evalContext := EvalContext{
		"env": "prod",
	}
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
	assert.True(t, store.IsEnabled("new_ui", evalContext))
}

func TestNewStoreFromURL_Errors(t *testing.T) {
	type args struct {
		statusCode int
	}
	tests := []struct {
		name     string
		args     args
		expected error
	}{
		{
			name: "304",
			args: args{
				statusCode: http.StatusNotModified,
			},
			expected: NotModifiedError,
		},
		{
			name: "500",
			args: args{
				statusCode: http.StatusInternalServerError,
			},
			expected: errors.New("server error: 500 Internal Server Error"),
		},
		{
			name: "200 (but empty body)",
			args: args{
				statusCode: http.StatusOK,
			},
			expected: errors.New("parse JSON: unexpected end of JSON input"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.statusCode)
				_, _ = w.Write([]byte(""))
			}))
			defer srv.Close()

			_, err := NewStoreFromURL(srv.URL+"/flags", "test-token", SetLastModified(time.Now()))
			assert.Equal(t, err.Error(), tt.expected.Error())
		})
	}
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
