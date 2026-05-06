// Package patcher applies a set of patch operations (set, unset, rename)
// to an environment map, returning a new map and a summary of changes.
package patcher

import "fmt"

// OpType represents the kind of patch operation.
type OpType string

const (
	OpSet    OpType = "set"
	OpUnset  OpType = "unset"
	OpRename OpType = "rename"
)

// Op describes a single patch operation.
type Op struct {
	Type  OpType
	Key   string
	Value string // used by OpSet
	To    string // used by OpRename
}

// Result summarises what the patch did.
type Result struct {
	Applied  []string
	Skipped  []string
	Warnings []string
}

// Patch applies ops to env and returns a new map plus a Result.
// The original map is never modified.
func Patch(env map[string]string, ops []Op) (map[string]string, Result) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var res Result

	for _, op := range ops {
		switch op.Type {
		case OpSet:
			out[op.Key] = op.Value
			res.Applied = append(res.Applied, fmt.Sprintf("set %s", op.Key))

		case OpUnset:
			if _, ok := out[op.Key]; !ok {
				res.Skipped = append(res.Skipped, fmt.Sprintf("unset %s (not found)", op.Key))
				continue
			}
			delete(out, op.Key)
			res.Applied = append(res.Applied, fmt.Sprintf("unset %s", op.Key))

		case OpRename:
			if _, ok := out[op.Key]; !ok {
				res.Skipped = append(res.Skipped, fmt.Sprintf("rename %s (not found)", op.Key))
				continue
			}
			if _, exists := out[op.To]; exists {
				res.Warnings = append(res.Warnings, fmt.Sprintf("rename %s->%s overwrites existing key", op.Key, op.To))
			}
			out[op.To] = out[op.Key]
			delete(out, op.Key)
			res.Applied = append(res.Applied, fmt.Sprintf("rename %s->%s", op.Key, op.To))

		default:
			res.Warnings = append(res.Warnings, fmt.Sprintf("unknown op type: %s", op.Type))
		}
	}

	return out, res
}
