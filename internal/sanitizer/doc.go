// Package sanitizer cleans environment variable values by removing or
// normalizing unsafe, invisible, or malformed characters.
//
// It supports stripping ASCII control characters, trimming whitespace,
// collapsing internal whitespace runs, and removing non-printable Unicode
// code points. All operations produce a new map, leaving the original
// untouched, and report which keys were affected.
package sanitizer
