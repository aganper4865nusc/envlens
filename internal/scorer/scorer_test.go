package scorer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/auditor"
	"github.com/yourorg/envlens/internal/linter"
	"github.com/yourorg/envlens/internal/scorer"
	"github.com/yourorg/envlens/internal/validator"
)

func TestScore_PerfectScore(t *testing.T) {
	r := scorer.Score(nil, nil, nil)
	if r.Score != 100 {
		t.Fatalf("expected 100, got %d", r.Score)
	}
	if r.Grade != scorer.GradeA {
		t.Fatalf("expected A, got %s", r.Grade)
	}
}

func TestScore_AuditHighSeverityDeducts15(t *testing.T) {
	issues := []auditor.Issue{{Message: "plain secret", Severity: auditor.SeverityHigh}}
	r := scorer.Score(issues, nil, nil)
	if r.Score != 85 {
		t.Fatalf("expected 85, got %d", r.Score)
	}
	if r.Grade != scorer.GradeB {
		t.Fatalf("expected B, got %s", r.Grade)
	}
}

func TestScore_LintErrorDeducts10(t *testing.T) {
	issues := []linter.Issue{{Message: "invalid char", Level: linter.LevelError}}
	r := scorer.Score(nil, issues, nil)
	if r.Score != 90 {
		t.Fatalf("expected 90, got %d", r.Score)
	}
}

func TestScore_ValidationIssueDeducts8(t *testing.T) {
	issues := []validator.Issue{{Message: "missing required key"}}
	r := scorer.Score(nil, nil, issues)
	if r.Score != 92 {
		t.Fatalf("expected 92, got %d", r.Score)
	}
}

func TestScore_NeverBelowZero(t *testing.T) {
	var auditIssues []auditor.Issue
	for i := 0; i < 20; i++ {
		auditIssues = append(auditIssues, auditor.Issue{Message: "issue", Severity: auditor.SeverityHigh})
	}
	r := scorer.Score(auditIssues, nil, nil)
	if r.Score != 0 {
		t.Fatalf("expected 0, got %d", r.Score)
	}
	if r.Grade != scorer.GradeF {
		t.Fatalf("expected F, got %s", r.Grade)
	}
}

func TestScore_DeductionsRecorded(t *testing.T) {
	issues := []linter.Issue{{Message: "trailing space", Level: linter.LevelWarn}}
	r := scorer.Score(nil, issues, nil)
	if len(r.Deductions) != 1 {
		t.Fatalf("expected 1 deduction, got %d", len(r.Deductions))
	}
	if r.Deductions[0].Points != 3 {
		t.Fatalf("expected 3 points deducted, got %d", r.Deductions[0].Points)
	}
}

func TestScore_GradeD(t *testing.T) {
	var vi []validator.Issue
	for i := 0; i < 7; i++ {
		vi = append(vi, validator.Issue{Message: "missing"})
	}
	r := scorer.Score(nil, nil, vi)
	// 100 - 7*8 = 44
	if r.Score != 44 {
		t.Fatalf("expected 44, got %d", r.Score)
	}
	if r.Grade != scorer.GradeD {
		t.Fatalf("expected D, got %s", r.Grade)
	}
}
