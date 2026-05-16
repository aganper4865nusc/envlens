package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func makeAnnotations() []differ.Annotation {
	return []differ.Annotation{
		{Key: "API_KEY", Status: "added", NewValue: "abc", Tags: []string{"added"}},
		{Key: "DB_PASS", Status: "removed", OldValue: "secret", Tags: []string{"removed"}},
		{Key: "PORT", Status: "changed", OldValue: "8080", NewValue: "9090", Tags: []string{"changed"}},
		{Key: "HOST", Status: "unchanged", OldValue: "localhost", NewValue: "localhost", Tags: []string{"unchanged"}},
	}
}

func TestWriteAnnotated_TextFormat_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteAnnotated(&buf, makeAnnotations(), "text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, key := range []string{"API_KEY", "DB_PASS", "PORT", "HOST"} {
		if !strings.Contains(out, key) {
			t.Errorf("output missing key %q", key)
		}
	}
}

func TestWriteAnnotated_TextFormat_ShowsSymbols(t *testing.T) {
	var buf bytes.Buffer
	WriteAnnotated(&buf, makeAnnotations(), "text")
	out := buf.String()

	if !strings.Contains(out, "[+]") {
		t.Error("expected [+] for added")
	}
	if !strings.Contains(out, "[-]") {
		t.Error("expected [-] for removed")
	}
	if !strings.Contains(out, "[~]") {
		t.Error("expected [~] for changed")
	}
	if !strings.Contains(out, "[=]") {
		t.Error("expected [=] for unchanged")
	}
}

func TestWriteAnnotated_TextFormat_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteAnnotated(&buf, []differ.Annotation{}, "text")
	out := buf.String()
	if !strings.Contains(out, "No diff annotations.") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestWriteAnnotated_JSONFormat_Structure(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteAnnotated(&buf, makeAnnotations(), "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 4 {
		t.Errorf("expected 4 records, got %d", len(records))
	}
	for _, r := range records {
		if _, ok := r["key"]; !ok {
			t.Error("missing 'key' field in JSON record")
		}
		if _, ok := r["status"]; !ok {
			t.Error("missing 'status' field in JSON record")
		}
		if _, ok := r["tags"]; !ok {
			t.Error("missing 'tags' field in JSON record")
		}
	}
}

func TestWriteAnnotated_TextFormat_ShowsTags(t *testing.T) {
	var buf bytes.Buffer
	anns := []differ.Annotation{
		{Key: "X", Status: "added", NewValue: "v", Tags: []string{"added", "env:prod"}},
	}
	WriteAnnotated(&buf, anns, "text")
	out := buf.String()
	if !strings.Contains(out, "env:prod") {
		t.Errorf("expected tag env:prod in output, got: %q", out)
	}
}
