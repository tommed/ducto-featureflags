package cli

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-featureflags/test"
	"testing"
)

func TestRunRoot_ListFlags(t *testing.T) {
	boolVars := test.Encode(test.BoolVariants(), "json")
	flags := fmt.Sprintf(`{
		"foo": { "variants": %s, "defaultVariant": "yes" },
		"bar": { "variants": %s, "defaultVariant": "no" }
	}`, boolVars, boolVars)
	path := writeTempFlags(t, flags)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	code := RunRoot([]string{"-file", path, "-list"}, stdout, stderr)

	assert.Equal(t, 0, code)
	output := stdout.String()
	assert.Contains(t, output, `"foo"`)
	assert.Contains(t, output, `"defaultVariant": "yes"`)
}
