package cli

import (
	"strings"

	"github.com/tommed/ducto-featureflags/sdk"
)

// For repeated --ctx key=value
type arrayFlags []string

func (a *arrayFlags) String() string { return strings.Join(*a, ",") }
func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

// Converts --ctx key=value into EvalContext
func parseContext(flags []string) sdk.EvalContext {
	ctx := sdk.EvalContext{}
	for _, f := range flags {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) == 2 {
			ctx[parts[0]] = parts[1]
		}
	}
	return ctx
}
