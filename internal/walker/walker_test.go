package walker_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlens/internal/walker"
)

func makeTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	write := func(rel string) {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("KEY=val\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(".env")
	write("staging/.env")
	write("staging/.env.local")
	write("production/.env")
	write("production/nested/deep/.env")
	write(".hidden/.env") // inside hidden dir
	write("notes.txt")    // should not match
	return root
}

func TestWalk_FindsDefaultPatterns(t *testing.T) {
	root := makeTree(t)
	results, err := walker.Walk(root, walker.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// expect: root/.env, staging/.env, staging/.env.local, production/.env,
	// production/nested/deep/.env, .hidden/.env
	if len(results) != 6 {
		t.Errorf("expected 6 results, got %d", len(results))
	}
}

func TestWalk_SkipHidden(t *testing.T) {
	root := makeTree(t)
	results, err := walker.Walk(root, walker.Options{SkipHidden: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		rel, _ := filepath.Rel(root, r.Path)
		if len(rel) > 0 && rel[0] == '.' && rel != ".env" {
			t.Errorf("hidden file should be skipped: %s", r.Path)
		}
	}
}

func TestWalk_MaxDepth(t *testing.T) {
	root := makeTree(t)
	results, err := walker.Walk(root, walker.Options{MaxDepth: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Depth > 1 {
			t.Errorf("depth %d exceeds max 1: %s", r.Depth, r.Path)
		}
	}
}

func TestWalk_CustomPattern(t *testing.T) {
	root := makeTree(t)
	results, err := walker.Walk(root, walker.Options{Patterns: []string{".env.local"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestWalk_SortedOutput(t *testing.T) {
	root := makeTree(t)
	results, err := walker.Walk(root, walker.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(results); i++ {
		if results[i].Path < results[i-1].Path {
			t.Errorf("results not sorted at index %d: %s < %s", i, results[i].Path, results[i-1].Path)
		}
	}
}

func TestWalk_EmptyDir(t *testing.T) {
	root := t.TempDir()
	results, err := walker.Walk(root, walker.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
