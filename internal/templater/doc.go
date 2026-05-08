// Package templater renders environment variable maps into arbitrary text
// templates using Go's text/template engine.
//
// Templates receive the env map as their data context, so any key can be
// referenced as .KEY_NAME. The package also exposes helper functions
// (default, upper, lower, warn) and supports configurable delimiters and
// missing-key policies via Options.
//
// Basic usage:
//
//	result, err := templater.Render("Hello, {{.NAME}}!", envMap, nil)
//
// With custom options:
//
//	opts := &templater.Options{
//		LeftDelim:  "[[",
//		RightDelim: "]]",
//		MissingKey: "error",
//	}
//	result, err := templater.Render("Hello, [[.NAME]]!", envMap, opts)
package templater
