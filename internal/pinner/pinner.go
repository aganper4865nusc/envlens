// Package pinner provides functionality to pin (lock) environment variable
// values to a snapshot, detecting any drift from the pinned state.
package pinner

import (
	"fmt"
	"sort"
)

// PinResult describes the outcome of comparing a live env against a pinned set.
type PinResult struct {
	Key      string
	Status   string // "ok", "drifted", "missing", "extra"
	Pinned   string
	Actual   string
}

// Options controls which keys to pin and how to handle extras.
type Options struct {
	// Keys is the explicit set of keys to pin. If empty, all pinned keys are checked.
	Keys []string
	// AllowExtra skips reporting keys present in live but absent from pinned.
	AllowExtra bool
}

// Pin compares a live environment map against a pinned (reference) map and
// returns a sorted slice of PinResult entries describing any drift.
func Pin(pinned, live map[string]string, opts Options) []PinResult {
	var results []PinResult

	keysToCheck := opts.Keys
	if len(keysToCheck) == 0 {
		for k := range pinned {
			keysToCheck = append(keysToCheck, k)
		}
	}
	sort.Strings(keysToCheck)

	checked := make(map[string]bool)
	for _, k := range keysToCheck {
		checked[k] = true
		pinnedVal, inPinned := pinned[k]
		actualVal, inLive := live[k]

		switch {
		case !inPinned:
			// Key requested but not in pinned set — skip silently
			continue
		case !inLive:
			results = append(results, PinResult{Key: k, Status: "missing", Pinned: pinnedVal, Actual: ""})
		case actualVal != pinnedVal:
			results = append(results, PinResult{Key: k, Status: "drifted", Pinned: pinnedVal, Actual: actualVal})
		default:
			results = append(results, PinResult{Key: k, Status: "ok", Pinned: pinnedVal, Actual: actualVal})
		}
	}

	if !opts.AllowExtra {
		var extraKeys []string
		for k := range live {
			if _, wasPinned := pinned[k]; !wasPinned {
				extraKeys = append(extraKeys, k)
			}
		}
		sort.Strings(extraKeys)
		for _, k := range extraKeys {
			results = append(results, PinResult{Key: k, Status: "extra", Pinned: "", Actual: live[k]})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return results
}

// HasDrift returns true if any result has a status other than "ok".
func HasDrift(results []PinResult) bool {
	for _, r := range results {
		if r.Status != "ok" {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary line for a PinResult.
func Summary(r PinResult) string {
	switch r.Status {
	case "ok":
		return fmt.Sprintf("[ok]      %s", r.Key)
	case "drifted":
		return fmt.Sprintf("[drifted] %s: pinned=%q actual=%q", r.Key, r.Pinned, r.Actual)
	case "missing":
		return fmt.Sprintf("[missing] %s: expected %q but key not present", r.Key, r.Pinned)
	case "extra":
		return fmt.Sprintf("[extra]   %s: not in pinned set (value=%q)", r.Key, r.Actual)
	default:
		return fmt.Sprintf("[unknown] %s", r.Key)
	}
}
