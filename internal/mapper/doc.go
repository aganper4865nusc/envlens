// Package mapper applies key-renaming and value-substitution rules to
// environment variable maps.
//
// Use [Map] with [Options] to rename keys via [KeyRule] entries and substitute
// specific values via [ValueRule] entries. Set DropUnmapped to true to produce
// a filtered map containing only keys referenced by a KeyRule.
//
// The returned [Result] records which keys were renamed and which values were
// substituted, enabling downstream audit or reporting.
package mapper
