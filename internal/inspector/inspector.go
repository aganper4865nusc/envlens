// Package inspector provides a high-level summary of a single env map,
// combining profiling, linting, auditing, and scoring into one Report.
package inspector

import (
	"github.com/yourusername/envlens/internal/auditor"
	"github.com/yourusername/envlens/internal/linter"
	"github.com/yourusername/envlens/internal/profiler"
	"github.com/yourusername/envlens/internal/scorer"
)

// Report holds the combined inspection result for an env map.
type Report struct {
	Source  string
	Profile profiler.Profile
	Lint    []linter.Issue
	Audit   []auditor.Issue
	Score   scorer.Result
}

// Inspect runs all analysis passes on env and returns a Report.
func Inspect(source string, env map[string]string) Report {
	prof := profiler.Analyze(env)

	lintIssues := linter.Lint(env)

	auditIssues := auditor.Audit(env)

	sc := scorer.Score(scorer.Input{
		AuditIssues:      auditIssues,
		LintIssues:       lintIssues,
		ValidationIssues: nil,
	})

	return Report{
		Source:  source,
		Profile: prof,
		Lint:    lintIssues,
		Audit:   auditIssues,
		Score:   sc,
	}
}

// HasIssues returns true if the report contains any lint or audit findings.
func HasIssues(r Report) bool {
	return len(r.Lint) > 0 || len(r.Audit) > 0
}
