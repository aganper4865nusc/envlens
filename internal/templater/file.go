package templater

import (
	"fmt"
	"os"
)

// RenderFile reads a template from srcPath, renders it against env using opts,
// and writes the result to dstPath. If dstPath is "-" the output is written to
// stdout. Intermediate directories are NOT created automatically.
func RenderFile(srcPath, dstPath string, env map[string]string, opts Options) (Result, error) {
	raw, err := os.ReadFile(srcPath)
	if err != nil {
		return Result{}, fmt.Errorf("templater: read %q: %w", srcPath, err)
	}

	res, err := Render(string(raw), env, opts)
	if err != nil {
		return Result{}, err
	}

	if dstPath == "-" {
		_, err = os.Stdout.WriteString(res.Output)
		if err != nil {
			return Result{}, fmt.Errorf("templater: write stdout: %w", err)
		}
		return res, nil
	}

	if err := os.WriteFile(dstPath, []byte(res.Output), 0o644); err != nil {
		return Result{}, fmt.Errorf("templater: write %q: %w", dstPath, err)
	}

	return res, nil
}
