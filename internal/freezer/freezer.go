// Package freezer provides functionality to lock an environment map
// to a read-only snapshot, preventing accidental mutation during audits
// or comparisons. A frozen env can be thawed back into a mutable map.
package freezer

import (
	"errors"
	"sort"
)

// ErrFrozen is returned when a mutation is attempted on a frozen environment.
var ErrFrozen = errors.New("environment is frozen and cannot be modified")

// Frozen holds an immutable copy of an environment map.
type Frozen struct {
	keys   []string
	values map[string]string
}

// Freeze creates a Frozen snapshot from the given environment map.
// The original map is not modified.
func Freeze(env map[string]string) *Frozen {
	keys := make([]string, 0, len(env))
	values := make(map[string]string, len(env))
	for k, v := range env {
		keys = append(keys, k)
		values[k] = v
	}
	sort.Strings(keys)
	return &Frozen{keys: keys, values: values}
}

// Get returns the value for the given key and whether it was found.
func (f *Frozen) Get(key string) (string, bool) {
	v, ok := f.values[key]
	return v, ok
}

// Keys returns a sorted slice of all keys in the frozen environment.
func (f *Frozen) Keys() []string {
	out := make([]string, len(f.keys))
	copy(out, f.keys)
	return out
}

// Len returns the number of keys in the frozen environment.
func (f *Frozen) Len() int {
	return len(f.keys)
}

// Thaw returns a mutable copy of the frozen environment map.
func (f *Frozen) Thaw() map[string]string {
	out := make(map[string]string, len(f.values))
	for k, v := range f.values {
		out[k] = v
	}
	return out
}

// Has reports whether the given key exists in the frozen environment.
func (f *Frozen) Has(key string) bool {
	_, ok := f.values[key]
	return ok
}
