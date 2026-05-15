package inverter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/inverter"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestInvert_BasicSwap(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	res := inverter.Invert(env, inverter.DefaultOptions())

	if res.Inverted["bar"] != "FOO" {
		t.Errorf("expected bar→FOO, got %q", res.Inverted["bar"])
	}
	if res.Inverted["qux"] != "BAZ" {
		t.Errorf("expected qux→BAZ, got %q", res.Inverted["qux"])
	}
	if len(res.Collisions) != 0 {
		t.Errorf("expected no collisions, got %v", res.Collisions)
	}
}

func TestInvert_CollisionFirstPolicy(t *testing.T) {
	// "ALPHA" and "BETA" both have value "shared"; alphabetically ALPHA wins.
	env := makeEnv("ALPHA", "shared", "BETA", "shared")
	res := inverter.Invert(env, inverter.Options{OnCollision: "first"})

	if res.Inverted["shared"] != "ALPHA" {
		t.Errorf("expected ALPHA to win, got %q", res.Inverted["shared"])
	}
	if len(res.Collisions["shared"]) != 2 {
		t.Errorf("expected 2 collision entries, got %v", res.Collisions["shared"])
	}
}

func TestInvert_CollisionLastPolicy(t *testing.T) {
	env := makeEnv("ALPHA", "shared", "BETA", "shared")
	res := inverter.Invert(env, inverter.Options{OnCollision: "last"})

	if res.Inverted["shared"] != "BETA" {
		t.Errorf("expected BETA to win, got %q", res.Inverted["shared"])
	}
}

func TestInvert_CollisionSkipPolicy(t *testing.T) {
	env := makeEnv("ALPHA", "shared", "BETA", "shared", "GAMMA", "unique")
	res := inverter.Invert(env, inverter.Options{OnCollision: "skip"})

	if _, ok := res.Inverted["shared"]; ok {
		t.Error("expected 'shared' to be skipped, but it is present")
	}
	if res.Inverted["unique"] != "GAMMA" {
		t.Errorf("expected unique→GAMMA, got %q", res.Inverted["unique"])
	}
}

func TestInvert_EmptyMap(t *testing.T) {
	res := inverter.Invert(map[string]string{}, inverter.DefaultOptions())
	if len(res.Inverted) != 0 {
		t.Errorf("expected empty inverted map, got %v", res.Inverted)
	}
}

func TestInvert_DefaultCollisionIsFirst(t *testing.T) {
	env := makeEnv("A", "v", "B", "v")
	res := inverter.Invert(env, inverter.Options{}) // zero value
	if res.Inverted["v"] != "A" {
		t.Errorf("expected A (first) to win with zero options, got %q", res.Inverted["v"])
	}
}
