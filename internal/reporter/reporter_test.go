package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/auditor"
	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/validator"
)

func TestWriteText_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	r := Report{}
	if err := Write(&buf, r, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No issues found.") {
		t.Errorf("expected 'No issues found.' but got: %s", buf.String())
	}
}

func TestWriteText_DiffAdded(t *testing.T) {
	var buf bytes.Buffer
	r := Report{
		DiffResults: []differ.DiffEntry{
			{Key: "NEW_KEY", Status: differ.StatusAdded, NewValue: "hello"},
		},
	}
	if err := Write(&buf, r, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY=hello") {
		t.Errorf("expected added key in output, got: %s", out)
	}
}

func TestWriteText_ValidationIssue(t *testing.T) {
	var buf bytes.Buffer
	r := Report{
		ValidateResults: []validator.ValidationIssue{
			{Key: "DB_URL", Severity: "error", Message: "required key is missing"},
		},
	}
	if err := Write(&buf, r, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[ERROR] DB_URL") {
		t.Errorf("expected validation issue in output, got: %s", out)
	}
}

func TestWriteText_AuditIssue(t *testing.T) {
	var buf bytes.Buffer
	r := Report{
		AuditResults: []auditor.AuditIssue{
			{Key: "SECRET", Severity: "warn", Message: "plain-text secret detected"},
		},
	}
	if err := Write(&buf, r, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[WARN] SECRET") {
		t.Errorf("expected audit issue in output, got: %s", out)
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	r := Report{
		DiffResults: []differ.DiffEntry{
			{Key: "FOO", Status: differ.StatusRemoved, OldValue: "bar"},
		},
	}
	if err := Write(&buf, r, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{`"diff"`, `"validation"`, `"audit"`, `"FOO"`} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output, got: %s", want, out)
		}
	}
}

func TestWriteJSON_EmptyReport(t *testing.T) {
	var buf bytes.Buffer
	r := Report{}
	if err := Write(&buf, r, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[]") {
		t.Errorf("expected empty arrays in JSON output, got: %s", out)
	}
}
