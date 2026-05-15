// Package pruner provides functionality to remove environment variable
// entries from a map based on configurable criteria.
//
// Entries can be pruned by:
//   - Exact key name (Keys)
//   - Key prefix (Prefixes)
//   - Regular expression match on the key (Patterns)
//   - Empty value (DropEmpty)
//
// All criteria are evaluated independently; a key is removed if any
// single criterion matches. The original map is never mutated.
package pruner
