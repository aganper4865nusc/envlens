// Package templater renders env maps into text templates using Go's text/template engine.
package templater

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Options controls template rendering behaviour.
type Options struct {
	// LeftDelim and RightDelim override the default {{ }} delimiters.
	LeftDelim  string
	RightDelim string
	// MissingKey controls behaviour when a key is absent: "zero", "error", or "default" (empty string).
	MissingKey string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		LeftDelim:  "{{",
		RightDelim: "}}",
		MissingKey: "error",
	}
}

// Result holds the rendered output and any warnings produced during rendering.
type Result struct {
	Output   string
	Warnings []string
}

// Render applies tmplSrc (a Go text/template string) against env and returns a Result.
func Render(tmplSrc string, env map[string]string, opts Options) (Result, error) {
	if opts.LeftDelim == "" {
		opts.LeftDelim = "{{"
	}
	if opts.RightDelim == "" {
		opts.RightDelim = "}}"
	}

	missingKey := opts.MissingKey
	if missingKey == "" {
		missingKey = "error"
	}

	var warnings []string

	funcMap := template.FuncMap{
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"warn": func(msg string) string {
			warnings = append(warnings, msg)
			return ""
		},
	}

	t, err := template.New("envlens").
		Delims(opts.LeftDelim, opts.RightDelim).
		Option("missingkey=" + missingKey).
		Funcs(funcMap).
		Parse(tmplSrc)
	if err != nil {
		return Result{}, fmt.Errorf("templater: parse error: %w", err)
	}

	data := make(map[string]string, len(env))
	for k, v := range env {
		data[k] = v
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return Result{}, fmt.Errorf("templater: render error: %w", err)
	}

	return Result{Output: buf.String(), Warnings: warnings}, nil
}
