package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/tommed/ducto-featureflags/sdk"
)

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

	store, err := sdk.NewFileWatchingStore(file)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load flags: %v\n", err)
		return 1
	}
	defer store.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/flags", func(w http.ResponseWriter, r *http.Request) {
		if token != "" {
			auth := r.Header.Get("Authorization")
			expected := "Bearer " + token
			if auth != expected {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		key := r.URL.Query().Get("key")
		if key != "" {
			// Convert query params to EvalContext
			ctx := sdk.EvalContext{}
			for k, v := range r.URL.Query() {
				if len(v) > 0 {
					ctx[k] = v[0]
				}
			}
			// Now fetch the flag value
			result := store.IsEnabled(key, ctx)
			_ = json.NewEncoder(w).Encode(map[string]bool{"enabled": result})
			return
		}

		// Just list all flags
		_ = json.NewEncoder(w).Encode(store.AllFlags())
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Fprintf(stdout, "Listening on %s...\n", addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintf(stderr, "server failed: %v", err)
		return 1
	}

	return 0
}
