package cli

import (
	"io"
)

func RunRoot(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return Run(args, stdout, stderr)
	}

	switch args[0] {
	case "serve":
		return Serve(args[1:], stdout, stderr)
	default:
		return Run(args, stdout, stderr)
	}
}
