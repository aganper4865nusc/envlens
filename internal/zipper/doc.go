// Package zipper provides a key-aligned merge of two environment variable maps
// via a user-supplied combiner function.
//
// Unlike merger (which performs a last-writer-wins union) or differ (which
// reports changes), zipper hands both values for every shared key to a
// CombineFunc so the caller can implement arbitrary merge strategies such as
// preferring one side, concatenating, or selecting the longer value.
//
// Keys present in only one map are included or excluded based on Options.
package zipper
