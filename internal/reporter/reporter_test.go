package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/auditor"
	"github.com/yourorg/envlens/internal/differ"
	"github.com/yourorg/envlens/internal/reporter"
	"github.com/yourorg/envlens/internal/resolver"
	"github.com/yourorg/envlens/internal/validator"
)

func TestWriteText_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	err := reporter.Write(&buf, reporter.Report{}, "text")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No issues") {
		t.Errorf("expected no-issues message, got: %s", buf.String())
	}
}

func TestWriteText_DiffAdded(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.Report{
		Diffs: []differ.Entry{{Key: "NEW_KEY", Status: differ.Added}},
	}
	_ = reporter.Write(&buf, r, "text")
	if !strings.Contains(buf.String(), "[diff]") {
		t.Errorf("expected diff line, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "NEW_KEY") {
		t.Errorf("expected key in output, got: %s", buf.String())
	}
}

func TestWriteText_ValidationIssue(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.Report{
		Validations: []validator.Issue{{Key: "DB_URL", Message: "required key missing"}},
	}
	_ = reporter.Write(&buf, r, "text")
	if !strings.Contains(buf.String(), "[validation]") {
		t.Errorf("expected validation line, got: %s", buf.String())
	}
}

func TestWriteText_AuditIssue(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.Report{
		Audits: []auditor.Finding{{Key: "SECRET", Message: "plain secret detected"}},
	}
	_ = reporter.Write(&buf, r, "text")
	if !strings.Contains(buf.String(), "[audit]") {
		t.Errorf("expected audit line, got: %s", buf.String())
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.Report{
		Diffs:       []differ.Entry{{Key: "X", Status: differ.Removed}},
		Validations: []validator.Issue{{Key: "Y", Message: "empty"}},
		Audits:      []auditor.Finding{{Key: "Z", Message: "plain"}},
	}
	_ = reporter.Write(&buf, r, "json")
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, field := range []string{"diffs", "validations", "audits"} {
		if _, ok := out[field]; !ok {
			t.Errorf("missing JSON field: %s", field)
		}
	}
}

func TestWriteJSON_WithResolutions(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.Report{
		Resolutions: []resolver.Resolution{
			{Key: "PORT", Value: "8080", Source: "file"},
			{Key: "HOST", Value: "localhost", Source: "default"},
		},
	}
	_ = reporter.Write(&buf, r, "json")
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	res, ok := out["resolutions"].([]interface{})
	if !ok || len(res) != 2 {
		t.Errorf("expected 2 resolutions, got: %v", out["resolutions"])
	}
}
