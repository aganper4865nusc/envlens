package inverter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/inverter"
)

// TestInvert_RoundTrip verifies that inverting twice returns the original map
// when there are no collisions.
func TestInvert_RoundTrip(t *testing.T) {
	original := makeEnv(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"APP_ENV", "production",
	)

	first := inverter.Invert(original, inverter.DefaultOptions())
	if len(first.Collisions) != 0 {
		t.Fatalf("unexpected collisions in first pass: %v", first.Collisions)
	}

	second := inverter.Invert(first.Inverted, inverter.DefaultOptions())
	if len(second.Collisions) != 0 {
		t.Fatalf("unexpected collisions in second pass: %v", second.Collisions)
	}

	for k, v := range original {
		if second.Inverted[k] != v {
			t.Errorf("round-trip mismatch for key %q: want %q got %q", k, v, second.Inverted[k])
		}
	}
	if len(second.Inverted) != len(original) {
		t.Errorf("length mismatch after round-trip: want %d got %d", len(original), len(second.Inverted))
	}
}

// TestInvert_CollisionsNotInInverted ensures that keys recorded as collisions
// under the "skip" policy are absent from the inverted map.
func TestInvert_SkipCollisionsAbsentFromResult(t *testing.T) {
	env := makeEnv(
		"X", "dup",
		"Y", "dup",
		"Z", "dup",
		"SOLO", "unique",
	)
	res := inverter.Invert(env, inverter.Options{OnCollision: "skip"})

	if _, present := res.Inverted["dup"]; present {
		t.Error("'dup' should not appear in inverted map under skip policy")
	}
	if res.Inverted["unique"] != "SOLO" {
		t.Errorf("expected unique→SOLO, got %q", res.Inverted["unique"])
	}
	// All three colliding original keys should be recorded.
	if len(res.Collisions["dup"]) != 3 {
		t.Errorf("expected 3 collision entries for 'dup', got %d", len(res.Collisions["dup"]))
	}
}
