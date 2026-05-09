// Package digester computes deterministic hashes of environment variable maps,
// enabling fingerprinting and change detection across deployments.
package digester

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Result holds the digest output for an env map.
type Result struct {
	// Fingerprint is the hex-encoded SHA-256 hash of the sorted key=value pairs.
	Fingerprint string
	// KeyCount is the number of keys included in the digest.
	KeyCount int
	// Keys lists the sorted keys that were hashed.
	Keys []string
}

// Options controls how the digest is computed.
type Options struct {
	// ExcludeKeys is a set of keys to omit from the digest.
	ExcludeKeys []string
	// KeysOnly hashes only key names, ignoring values.
	KeysOnly bool
}

// Digest computes a deterministic fingerprint of the given env map.
func Digest(env map[string]string, opts Options) Result {
	excluded := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[k] = true
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		if !excluded[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		var entry string
		if opts.KeysOnly {
			entry = k
		} else {
			entry = fmt.Sprintf("%s=%s", k, env[k])
		}
		fmt.Fprintln(h, entry)
	}

	return Result{
		Fingerprint: hex.EncodeToString(h.Sum(nil)),
		KeyCount:    len(keys),
		Keys:        keys,
	}
}

// Equal returns true if two Results share the same fingerprint.
func Equal(a, b Result) bool {
	return strings.EqualFold(a.Fingerprint, b.Fingerprint)
}
