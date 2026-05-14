// Package censor provides utilities for censoring environment variable
// values based on configurable rules, replacing sensitive data with
// placeholder tokens suitable for logging or display.
package censor

import (
	"regexp"
	"strings"
)

// Options controls censor behaviour.
type Options struct {
	// Token is the replacement string for censored values. Defaults to "[CENSORED]".
	Token string
	// KeyPatterns is a list of regex patterns matched against keys (case-insensitive).
	KeyPatterns []string
	// ValuePatterns is a list of regex patterns; if a value matches any, it is censored.
	ValuePatterns []string
}

// Result holds the censored environment map and metadata.
type Result struct {
	Env      map[string]string
	Censored []string // keys whose values were censored
}

var defaultKeyPatterns = []string{
	`(?i)(secret|password|passwd|token|apikey|api_key|private_key|credential|auth)`,
}

// Censor returns a copy of env with sensitive values replaced by the token.
// Keys are matched against KeyPatterns; values are matched against ValuePatterns.
func Censor(env map[string]string, opts Options) Result {
	if opts.Token == "" {
		opts.Token = "[CENSORED]"
	}
	if len(opts.KeyPatterns) == 0 {
		opts.KeyPatterns = defaultKeyPatterns
	}

	keyRegs := compileAll(opts.KeyPatterns)
	valRegs := compileAll(opts.ValuePatterns)

	out := make(map[string]string, len(env))
	var censored []string

	for k, v := range env {
		if v == "" {
			out[k] = v
			continue
		}
		if matchesAny(keyRegs, k) || matchesAny(valRegs, v) {
			out[k] = opts.Token
			censored = append(censored, k)
		} else {
			out[k] = v
		}
	}

	sortStrings(censored)
	return Result{Env: out, Censored: censored}
}

func compileAll(patterns []string) []*regexp.Regexp {
	var regs []*regexp.Regexp
	for _, p := range patterns {
		if r, err := regexp.Compile(p); err == nil {
			regs = append(regs, r)
		}
	}
	return regs
}

func matchesAny(regs []*regexp.Regexp, s string) bool {
	for _, r := range regs {
		if r.MatchString(s) {
			return true
		}
	}
	return false
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && strings.ToLower(ss[j]) < strings.ToLower(ss[j-1]); j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
