package sdk

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-featureflags/test"
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
			"variants": ` + test.BoolVariantsJSON() + `,
			"defaultVariant": "yes"
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

	newUIFlag, ok := store.Get("new_ui")
	assert.True(t, ok)
	result := newUIFlag.Evaluate(evalContext)

	assert.True(t, result.Value.(bool))
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
  variants:
    yes: true
    no: false
  defaultVariant: no
  rules:
  - if:
      env: prod
    variant: yes
`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		_, _ = w.Write([]byte(payload))
	}))
	defer srv.Close()

	ctx := context.Background()
	store, err := NewStoreFromURL(ctx, srv.URL+"/flags.yaml", "")
	assert.NoError(t, err)

	prodCtx := EvalContext{"env": "prod"}
	devCtx := EvalContext{"env": "dev"}

	prodEnabledFlag, _ := store.Get("canary")
	devEnabledFlag, _ := store.Get("canary")
	prodEnabled := prodEnabledFlag.Evaluate(prodCtx)
	devEnabled := devEnabledFlag.Evaluate(devCtx)
	assert.True(t, prodEnabled.Value.(bool), "prod")
	assert.False(t, devEnabled.Value.(bool), "dev")
}
