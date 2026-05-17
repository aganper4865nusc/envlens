// Package zipper provides a key-aligned merge of two environment variable maps
// via a user-supplied combiner function.
//
// Unlike merger (which performs a last-writer-wins union) or differ (which
// reports changes), zipper hands both values for every shared key to a
// CombineFunc so the caller can implement arbitrary merge strategies such as
// preferring one side, concatenating, or selecting the longer value.
//
// Keys present in only one map are included or excluded based on Options.
//
// # Basic usage
//
//	result, err := zipper.Zip(base, override, func(key, a, b string) (string, error) {
//		// prefer the override value when non-empty, otherwise keep base
//		if b != "" {
//			return b, nil
//		}
//		return a, nil
//	}, nil)
//
// # Options
//
// Pass a non-nil *Options to control how unmatched keys (keys present in only
// one of the two input maps) are handled:
//
//   - IncludeLeftOnly  – retain keys found only in the left/base map.
//   - IncludeRightOnly – retain keys found only in the right/override map.
package zipper
