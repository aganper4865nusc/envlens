// Package freezer provides a Frozen type that wraps an environment map
// in an immutable, read-only snapshot. Use Freeze to lock a map and
// Thaw to retrieve a mutable copy when modifications are needed again.
//
// Typical usage:
//
//	f := freezer.Freeze(env)
//	v, ok := f.Get("DB_HOST")
//	mutable := f.Thaw()
package freezer
