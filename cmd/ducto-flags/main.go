// File: ducto-featureflags/cmd/ducto-flags/main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tommed/ducto-featureflags/sdk"
)

func main() {
	// CLI flags
	file := flag.String("file", "flags.json", "Path to feature flag definition file")
	key := flag.String("key", "", "Feature flag key to check")
	printAll := flag.Bool("list", false, "Print all loaded flags")

	flag.Parse()

	if *file == "" {
		log.Fatal("missing required flag: -file")
	}

	// Load the store
	store, err := sdk.NewStoreFromFile(*file)
	if err != nil {
		log.Fatalf("failed to load flags: %v", err)
	}

	if *printAll {
		flags := store.AllFlags()
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(flags)
		return
	}

	if *key == "" {
		log.Fatal("missing required flag: -key (or use -list to dump all flags)")
	}

	enabled := store.IsEnabled(*key)
	fmt.Printf("Flag %q is %v\n", *key, enabled)
}
