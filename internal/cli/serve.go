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

	fs.StringVar(&file, "file", "flags.json", "Path to feature flag definition file")
	fs.StringVar(&addr, "addr", ":8080", "Listen address")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "failed to parse serve flags: %v", err)
		return 1
	}

	store, err := sdk.NewWatchingStore(file)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load flags: %v", err)
		return 1
	}
	defer store.Close()

	http.HandleFunc("/api/flags", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")

		// Convert query params to EvalContext
		ctx := sdk.EvalContext{}
		for k, v := range r.URL.Query() {
			if len(v) > 0 {
				ctx[k] = v[0]
			}
		}

		// Return single flag value
		if key != "" {
			result := store.IsEnabled(key, ctx)
			_ = json.NewEncoder(w).Encode(map[string]bool{"enabled": result})
			return
		}

		// Return all flags (raw definition)
		_ = json.NewEncoder(w).Encode(store.AllFlags())
	})

	fmt.Fprintf(stdout, "Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(stderr, "server failed: %v", err)
		return 1
	}

	return 0
}
