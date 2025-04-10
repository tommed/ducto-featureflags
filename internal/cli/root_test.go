package cli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunRoot_ListFlags(t *testing.T) {
	flags := `{
		"foo": { "enabled": true },
		"bar": { "enabled": false }
	}`
	path := writeTempFlags(t, flags)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	code := RunRoot([]string{"-file", path, "-list"}, stdout, stderr)

	assert.Equal(t, 0, code)
	output := stdout.String()
	assert.Contains(t, output, `"foo"`)
	assert.Contains(t, output, `"enabled": true`)
}
