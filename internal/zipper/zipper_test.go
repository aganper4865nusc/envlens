package zipper_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/zipper"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func preferLeft(l, r string) string {
	if l != "" {
		return l
	}
	return r
}

func TestZip_BothKeys(t *testing.T) {
	left := makeEnv("HOST", "localhost")
	right := makeEnv("HOST", "prod.example.com")
	results := zipper.Zip(left, right, zipper.DefaultOptions(), preferLeft)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Result != "localhost" {
		t.Errorf("expected localhost, got %s", results[0].Result)
	}
}

func TestZip_LeftOnly_Included(t *testing.T) {
	left := makeEnv("ONLY_LEFT", "yes")
	right := makeEnv[string, string]()
	opts := zipper.DefaultOptions()
	results := zipper.Zip(left, right, opts, preferLeft)
	if len(results) != 1 || !results[0].LeftOnly {
		t.Errorf("expected one left-only result")
	}
}

func TestZip_LeftOnly_Excluded(t *testing.T) {
	left := makeEnv("ONLY_LEFT", "yes")
	right := makeEnv[string, string]()
	opts := zipper.Options{IncludeLeftOnly: false, IncludeRightOnly: true}
	results := zipper.Zip(left, right, opts, preferLeft)
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestZip_RightOnly_Excluded(t *testing.T) {
	left := makeEnv[string, string]()
	right := makeEnv("ONLY_RIGHT", "val")
	opts := zipper.Options{IncludeLeftOnly: true, IncludeRightOnly: false}
	results := zipper.Zip(left, right, opts, preferLeft)
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestZip_SortedOutput(t *testing.T) {
	left := makeEnv("Z_KEY", "z", "A_KEY", "a", "M_KEY", "m")
	right := makeEnv("Z_KEY", "z2", "A_KEY", "a2", "M_KEY", "m2")
	results := zipper.Zip(left, right, zipper.DefaultOptions(), preferLeft)
	keys := make([]string, len(results))
	for i, r := range results {
		keys[i] = r.Key
	}
	joined := strings.Join(keys, ",")
	if joined != "A_KEY,M_KEY,Z_KEY" {
		t.Errorf("unexpected order: %s", joined)
	}
}

func TestToMap_Conversion(t *testing.T) {
	left := makeEnv("PORT", "8080", "HOST", "localhost")
	right := makeEnv("PORT", "9090", "HOST", "remote")
	results := zipper.Zip(left, right, zipper.DefaultOptions(), func(l, r string) string {
		return l + "|" + r
	})
	m := zipper.ToMap(results)
	if m["PORT"] != "8080|9090" {
		t.Errorf("unexpected PORT value: %s", m["PORT"])
	}
	if m["HOST"] != "localhost|remote" {
		t.Errorf("unexpected HOST value: %s", m["HOST"])
	}
}
