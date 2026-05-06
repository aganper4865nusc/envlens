package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envlens/internal/snapshot"
)

func TestTake_CopiesEnv(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.Take("test", env)
	env["FOO"] = "mutated"
	if s.Env["FOO"] != "bar" {
		t.Errorf("expected snapshot to be independent of original map")
	}
}

func TestTake_SetsSourceAndTimestamp(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.Take("staging", map[string]string{})
	if s.Source != "staging" {
		t.Errorf("expected source 'staging', got %q", s.Source)
	}
	if s.Timestamp.Before(before) {
		t.Errorf("expected timestamp to be set after test start")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	env := map[string]string{"KEY": "value", "OTHER": "123"}
	s := snapshot.Take("prod", env)

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(s, tmp); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Source != "prod" {
		t.Errorf("expected source 'prod', got %q", loaded.Source)
	}
	if loaded.Env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", loaded.Env["KEY"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(tmp, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := snapshot.Load(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
