package openfeature

import (
	"context"
	"fmt"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/tommed/ducto-featureflags/sdk"
	"time"
)

func ExampleDuctoProvider() {
	// Create the flag store (e.g., from a static file or dynamic source)
	store, _ := sdk.NewStoreFromBytesWithFormat([]byte(`{
		"my_flag": {
			"enabled": true,
			"defaultVariant": "on",
			"variants": {
				"on": "hello world"
			}
		}
	}`), "json")

	// Set Ducto as the global provider
	_ = openfeature.SetProvider(NewProvider(store))

	client := openfeature.NewClient("example-service")

	// Generate your eval context
	evalCtx := openfeature.NewEvaluationContext("test", map[string]interface{}{
		"test-group": "examples",
		"foo":        123,
	})

	// OpenFeature seems to need this, but not Ducto
	time.Sleep(10 * time.Millisecond)

	// Evaluate the flag
	val, err := client.StringValue(context.Background(), "my_flag", "fallback", evalCtx)
	fmt.Println(val, err)

	// Output:
	// hello world <nil>
}
