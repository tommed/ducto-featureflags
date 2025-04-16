package openfeature

import (
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/tommed/ducto-featureflags/sdk"
)

// convertFlattenedContext converts OpenFeature FlattenedContext to internal EvalContext
func convertFlattenedContext(fc openfeature.FlattenedContext) sdk.EvalContext {
	result := make(sdk.EvalContext)
	for k, v := range fc {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}
	return result
}
