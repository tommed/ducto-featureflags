package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/tommed/ducto-featureflags/sdk"
)

type ResolutionResponse struct {
	Variant string      `json:"variant"`
	Value   interface{} `json:"value"`
	Reason  string      `json:"reason"`
	Error   string      `json:"error,omitempty"`
}

//goland:noinspection GoUnhandledErrorResult
func Serve(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)

	var file string
	var addr string
	var token string

	fs.StringVar(&file, "file", "flags.json", "Path to feature flag definition file")
	fs.StringVar(&addr, "addr", ":8080", "Listen address")
	fs.StringVar(&token, "token", "", "Optional bearer token required to access the API")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "failed to parse serve flags: %v\n", err)
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	provider := sdk.NewFileProviderWithLog(file, stdout)
	store := sdk.NewDynamicStore(ctx, provider)
	err := store.Start()
	if err != nil {
		fmt.Fprintf(stderr, "failed to load flags: %v\n", err)
		return 1
	}

	mux := http.NewServeMux()
	var handler = func(encode func(w http.ResponseWriter, graph interface{})) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if token != "" {
				auth := r.Header.Get("Authorization")
				expected := "Bearer " + token
				if auth != expected {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
			}

			// Handle If-Modified-Since and 304s
			ifModified := r.Header.Get("If-Modified-Since")
			if ifModified != "" {
				if t, err := http.ParseTime(ifModified); err == nil {
					if !store.LastUpdated().IsZero() && store.LastUpdated().Before(t.Add(1*time.Second)) {
						w.WriteHeader(http.StatusNotModified)
						return
					}
				}
			}
			w.Header().Set("Last-Modified", store.LastUpdated().UTC().Format(http.TimeFormat))

			// Fetch the eval context from the query-string
			key := r.URL.Query().Get("key")
			if key != "" {
				// Convert query params to EvalContext
				ctx := sdk.EvalContext{}
				for k, v := range r.URL.Query() {
					if len(v) > 0 {
						ctx[k] = v[0]
					}
				}

				// Determine what to send back
				storeFlag, ok := store.Get(key)
				if !ok {
					encode(w, ResolutionResponse{
						Variant: "",
						Value:   false,
						Reason:  "ERROR",
						Error:   "flag not found",
					})
					return
				}

				variant, val, _, matched := storeFlag.Evaluate(ctx)
				resp := ResolutionResponse{
					Variant: variant,
					Value:   val,
					Reason:  "FALLBACK",
				}
				if matched {
					resp.Reason = "TARGETING_MATCH"
				}
				encode(w, resp)
				return
			}

			// Just list all flags
			encode(w, store.AllFlags())
		}
	}
	mux.HandleFunc("/api/flags", handler(handleJSON))
	mux.HandleFunc("/api/flags.yaml", handler(handleYAML))
	mux.HandleFunc("/api/flags.json", handler(handleJSON))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
		stop()
	}()

	fmt.Fprintf(stdout, "Listening on %s...\n", addr)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(stderr, "server failed: %v", err)
			return 1
		}
	}

	return 0
}

func handleJSON(w http.ResponseWriter, graph interface{}) {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	_ = e.Encode(graph)
}

func handleYAML(w http.ResponseWriter, graph interface{}) {
	_ = yaml.NewEncoder(w).Encode(graph)
}
