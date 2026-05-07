package pinner_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/pinner"
)

func TestPin_AllOk(t *testing.T) {
	pinned := map[string]string{"HOST": "localhost", "PORT": "8080"}
	live := map[string]string{"HOST": "localhost", "PORT": "8080"}

	results := pinner.Pin(pinned, live, pinner.Options{AllowExtra: true})
	for _, r := range results {
		if r.Status != "ok" {
			t.Errorf("expected ok for %s, got %s", r.Key, r.Status)
		}
	}
}

func TestPin_DriftedValue(t *testing.T) {
	pinned := map[string]string{"DB_HOST": "prod-db.internal"}
	live := map[string]string{"DB_HOST": "staging-db.internal"}

	results := pinner.Pin(pinned, live, pinner.Options{AllowExtra: true})
	if len(results) != 1 || results[0].Status != "drifted" {
		t.Fatalf("expected 1 drifted result, got %+v", results)
	}
	if results[0].Pinned != "prod-db.internal" {
		t.Errorf("unexpected pinned value: %s", results[0].Pinned)
	}
}

func TestPin_MissingKey(t *testing.T) {
	pinned := map[string]string{"SECRET_KEY": "abc123"}
	live := map[string]string{}

	results := pinner.Pin(pinned, live, pinner.Options{AllowExtra: true})
	if len(results) != 1 || results[0].Status != "missing" {
		t.Fatalf("expected missing result, got %+v", results)
	}
}

func TestPin_ExtraKey_Reported(t *testing.T) {
	pinned := map[string]string{"APP": "envlens"}
	live := map[string]string{"APP": "envlens", "EXTRA_VAR": "surprise"}

	results := pinner.Pin(pinned, live, pinner.Options{AllowExtra: false})
	statuses := map[string]string{}
	for _, r := range results {
		statuses[r.Key] = r.Status
	}
	if statuses["EXTRA_VAR"] != "extra" {
		t.Errorf("expected EXTRA_VAR to be extra, got %s", statuses["EXTRA_VAR"])
	}
}

func TestPin_ExtraKey_Allowed(t *testing.T) {
	pinned := map[string]string{"APP": "envlens"}
	live := map[string]string{"APP": "envlens", "EXTRA_VAR": "surprise"}

	results := pinner.Pin(pinned, live, pinner.Options{AllowExtra: true})
	for _, r := range results {
		if r.Key == "EXTRA_VAR" {
			t.Error("EXTRA_VAR should not appear when AllowExtra is true")
		}
	}
}

func TestPin_SpecificKeys(t *testing.T) {
	pinned := map[string]string{"A": "1", "B": "2", "C": "3"}
	live := map[string]string{"A": "1", "B": "changed", "C": "3"}

	results := pinner.Pin(pinned, live, pinner.Options{Keys: []string{"A", "C"}, AllowExtra: true})
	if len(results) != 2 {
		t.Fatalf("expected 2 results for keys A and C, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != "ok" {
			t.Errorf("key %s should be ok, got %s", r.Key, r.Status)
		}
	}
}

func TestHasDrift_True(t *testing.T) {
	results := []pinner.PinResult{{Key: "X", Status: "drifted"}}
	if !pinner.HasDrift(results) {
		t.Error("expected HasDrift to return true")
	}
}

func TestHasDrift_False(t *testing.T) {
	results := []pinner.PinResult{{Key: "X", Status: "ok"}}
	if pinner.HasDrift(results) {
		t.Error("expected HasDrift to return false")
	}
}

func TestSummary_Formats(t *testing.T) {
	cases := []struct {
		status string
		contains string
	}{
		{"ok", "[ok]"},
		{"drifted", "[drifted]"},
		{"missing", "[missing]"},
		{"extra", "[extra]"},
	}
	for _, tc := range cases {
		r := pinner.PinResult{Key: "K", Status: tc.status, Pinned: "p", Actual: "a"}
		s := pinner.Summary(r)
		if len(s) == 0 {
			t.Errorf("empty summary for status %s", tc.status)
		}
	}
}
