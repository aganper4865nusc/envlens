package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
)

func makeThreeWayResults() []differ.ThreeWayResult {
	return []differ.ThreeWayResult{
		{Key: "APP_ENV", Base: "dev", Left: "dev", Right: "dev", Conflict: false, Resolution: "dev"},
		{Key: "DB_HOST", Base: "localhost", Left: "db.prod", Right: "localhost", Conflict: false, Resolution: "db.prod"},
		{Key: "SECRET", Base: "old", Left: "left_val", Right: "right_val", Conflict: true, Resolution: ""},
	}
}

func TestWriteThreeWay_TextFormat_ContainsConflict(t *testing.T) {
	var buf bytes.Buffer
	err := WriteThreeWay(&buf, makeThreeWayResults(), "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[CONFLICT]") {
		t.Error("expected [CONFLICT] label in output")
	}
	if !strings.Contains(out, "SECRET") {
		t.Error("expected key SECRET in output")
	}
}

func TestWriteThreeWay_TextFormat_ShowsChangedKey(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteThreeWay(&buf, makeThreeWayResults(), "text")
	out := buf.String()
	if !strings.Contains(out, "[CHANGED]") {
		t.Error("expected [CHANGED] label for DB_HOST")
	}
}

func TestWriteThreeWay_TextFormat_ShowsOKKey(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteThreeWay(&buf, makeThreeWayResults(), "text")
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Error("expected [OK] label for APP_ENV")
	}
}

func TestWriteThreeWay_TextFormat_SummaryLine(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteThreeWay(&buf, makeThreeWayResults(), "text")
	out := buf.String()
	if !strings.Contains(out, "3 keys") {
		t.Error("expected total key count in summary")
	}
	if !strings.Contains(out, "1 conflict") {
		t.Error("expected conflict count in summary")
	}
}

func TestWriteThreeWay_JSONFormat_Structure(t *testing.T) {
	var buf bytes.Buffer
	err := WriteThreeWay(&buf, makeThreeWayResults(), "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out struct {
		TotalKeys int `json:"total_keys"`
		Conflicts int `json:"conflicts"`
		Entries   []struct {
			Key      string `json:"key"`
			Conflict bool   `json:"conflict"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", out.TotalKeys)
	}
	if out.Conflicts != 1 {
		t.Errorf("expected 1 conflict, got %d", out.Conflicts)
	}
	if out.Entries[2].Key != "SECRET" || !out.Entries[2].Conflict {
		t.Error("expected SECRET to be marked as conflict")
	}
}

func TestWriteThreeWay_DefaultFormat_IsText(t *testing.T) {
	var buf bytes.Buffer
	_ = WriteThreeWay(&buf, makeThreeWayResults(), "")
	out := buf.String()
	if !strings.Contains(out, "Three-Way Diff") {
		t.Error("expected text header for default format")
	}
}
