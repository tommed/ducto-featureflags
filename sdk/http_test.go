package sdk

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewStoreFromURL_JSON(t *testing.T) {
	evalContext := EvalContext{
		"env": "prod",
	}
	payload := `{
		"new_ui": {
			"enabled": true
		}
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()

	ctx := context.Background()
	store, err := NewStoreFromURL(ctx, srv.URL+"/flags", "test-token")
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
			name: "500",
			args: args{
				statusCode: http.StatusInternalServerError,
			},
			expected: errors.New("http error: 500 Internal Server Error"),
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

			ctx := context.Background()
			_, err := NewStoreFromURL(ctx, srv.URL+"/flags", "test-token")
			assert.Equal(t, tt.expected.Error(), err.Error())
		})
	}
}

func TestNewStoreFromURL_YAML(t *testing.T) {
	payload := `
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

	ctx := context.Background()
	store, err := NewStoreFromURL(ctx, srv.URL+"/flags.yaml", "")
	assert.NoError(t, err)

	prodEnabled := store.IsEnabled("canary", EvalContext{"env": "prod"})
	devEnabled := store.IsEnabled("canary", EvalContext{"env": "dev"})
	assert.True(t, prodEnabled, "prod")
	assert.False(t, devEnabled, "dev")
}
