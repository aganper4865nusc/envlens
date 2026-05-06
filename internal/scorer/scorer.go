// Package scorer computes a quality score for an environment variable file
// based on audit, lint, and validation results.
package scorer

import (
	"github.com/yourorg/envlens/internal/auditor"
	"github.com/yourorg/envlens/internal/linter"
	"github.com/yourorg/envlens/internal/validator"
)

// Grade represents a letter grade assigned to an env file.
type Grade string

const (
	GradeA Grade = "A"
	GradeB Grade = "B"
	GradeC Grade = "C"
	GradeD Grade = "D"
	GradeF Grade = "F"
)

// Result holds the numeric score and derived grade.
type Result struct {
	Score      int
	Grade      Grade
	Deductions []Deduction
}

// Deduction records why points were subtracted.
type Deduction struct {
	Reason string
	Points int
}

// Score computes a 0–100 quality score from audit, lint, and validation issues.
func Score(
	auditIssues []auditor.Issue,
	lintIssues []linter.Issue,
	validationIssues []validator.Issue,
) Result {
	score := 100
	var deductions []Deduction

	for _, issue := range auditIssues {
		pts := 5
		if issue.Severity == auditor.SeverityHigh {
			pts = 15
		}
		score -= pts
		deductions = append(deductions, Deduction{Reason: "audit: " + issue.Message, Points: pts})
	}

	for _, issue := range lintIssues {
		pts := 3
		if issue.Level == linter.LevelError {
			pts = 10
		}
		score -= pts
		deductions = append(deductions, Deduction{Reason: "lint: " + issue.Message, Points: pts})
	}

	for _, issue := range validationIssues {
		pts := 8
		score -= pts
		deductions = append(deductions, Deduction{Reason: "validation: " + issue.Message, Points: pts})
	}

	if score < 0 {
		score = 0
	}

	return Result{
		Score:      score,
		Grade:      toGrade(score),
		Deductions: deductions,
	}
}

func toGrade(score int) Grade {
	switch {
	case score >= 90:
		return GradeA
	case score >= 75:
		return GradeB
	case score >= 60:
		return GradeC
	case score >= 40:
		return GradeD
	default:
		return GradeF
	}
}
