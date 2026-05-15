// Package rotator provides utilities for rotating environment variable keys
// by applying a rename mapping with optional value transformation.
package rotator

import (
	"fmt"
	"strings"
)

// Op describes a single rotation operation.
type Op struct {
	FromKey   string
	ToKey     string
	Transform func(string) string // optional value transform
}

// Result holds the output of a Rotate call.
type Result struct {
	Env      map[string]string
	Rotated  []string // keys that were successfully rotated
	Skipped  []string // FromKeys not found in source
	Conflict []string // ToKeys that already existed (and were overwritten)
}

// Rotate applies the given ops to env, returning a new map and a Result.
// The original map is never modified.
func Rotate(env map[string]string, ops []Op) (Result, error) {
	out := copyMap(env)
	result := Result{Env: out}

	for _, op := range ops {
		if op.FromKey == "" || op.ToKey == "" {
			return Result{}, fmt.Errorf("rotator: op has empty FromKey or ToKey")
		}
		if strings.ContainsAny(op.ToKey, " \t") {
			return Result{}, fmt.Errorf("rotator: ToKey %q contains whitespace", op.ToKey)
		}

		val, ok := out[op.FromKey]
		if !ok {
			result.Skipped = append(result.Skipped, op.FromKey)
			continue
		}

		if _, exists := out[op.ToKey]; exists && op.ToKey != op.FromKey {
			result.Conflict = append(result.Conflict, op.ToKey)
		}

		if op.Transform != nil {
			val = op.Transform(val)
		}

		delete(out, op.FromKey)
		out[op.ToKey] = val
		result.Rotated = append(result.Rotated, op.FromKey)
	}

	return result, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
