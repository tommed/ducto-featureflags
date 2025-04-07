package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/tommed/ducto-featureflags/sdk"
)

//goland:noinspection GoUnhandledErrorResult
func Run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ducto-flags", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var file string
	var key string
	var printAll bool
	var ctxFlags arrayFlags

	fs.StringVar(&file, "file", "flags.json", "Path to feature flag definition file")
	fs.StringVar(&key, "key", "", "Feature flag key to check")
	fs.BoolVar(&printAll, "list", false, "Print all loaded flags")
	fs.Var(&ctxFlags, "ctx", "Context key=value pair (can be used multiple times)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "failed to parse args: %v\n", err)
		return 1
	}

	if file == "" {
		fmt.Fprintf(stderr, "missing required flag: -file")
		return 1
	}

	if key == "" && !printAll {
		fmt.Fprintf(stderr, "must provide either -key or -list")
		return 1
	}

	store, err := sdk.NewStoreFromFile(file)
	if err != nil {
		fmt.Fprintf(stderr, "failed to load flags: %v", err)
		return 1
	}

	if printAll {
		dumpAllFlags(stdout, store)
		return 0
	}

	ctx := parseContext(ctxFlags)
	result := store.IsEnabled(key, ctx)
	fmt.Fprintf(stdout, "Flag %q is %v\n", key, result)
	return 0
}

func dumpAllFlags(stdout io.Writer, store *sdk.Store) {
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(store.AllFlags())
}
