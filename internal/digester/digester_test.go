package digester_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/digester"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestDigest_DeterministicOutput(t *testing.T) {
	env := makeEnv("B", "2", "A", "1", "C", "3")
	r1 := digester.Digest(env, digester.Options{})
	r2 := digester.Digest(env, digester.Options{})
	if r1.Fingerprint != r2.Fingerprint {
		t.Errorf("expected same fingerprint, got %q vs %q", r1.Fingerprint, r2.Fingerprint)
	}
}

func TestDigest_DifferentValues_DifferentHash(t *testing.T) {
	a := digester.Digest(makeEnv("KEY", "val1"), digester.Options{})
	b := digester.Digest(makeEnv("KEY", "val2"), digester.Options{})
	if a.Fingerprint == b.Fingerprint {
		t.Error("expected different fingerprints for different values")
	}
}

func TestDigest_KeyCount(t *testing.T) {
	env := makeEnv("A", "1", "B", "2", "C", "3")
	r := digester.Digest(env, digester.Options{})
	if r.KeyCount != 3 {
		t.Errorf("expected KeyCount=3, got %d", r.KeyCount)
	}
}

func TestDigest_ExcludeKeys(t *testing.T) {
	env := makeEnv("A", "1", "SECRET", "topsecret", "B", "2")
	opts := digester.Options{ExcludeKeys: []string{"SECRET"}}
	r := digester.Digest(env, opts)
	if r.KeyCount != 2 {
		t.Errorf("expected KeyCount=2 after exclusion, got %d", r.KeyCount)
	}
	for _, k := range r.Keys {
		if k == "SECRET" {
			t.Error("SECRET should have been excluded from digest")
		}
	}
}

func TestDigest_KeysOnly_IgnoresValues(t *testing.T) {
	env1 := makeEnv("A", "hello", "B", "world")
	env2 := makeEnv("A", "different", "B", "values")
	opts := digester.Options{KeysOnly: true}
	r1 := digester.Digest(env1, opts)
	r2 := digester.Digest(env2, opts)
	if r1.Fingerprint != r2.Fingerprint {
		t.Error("KeysOnly: same keys with different values should produce same fingerprint")
	}
}

func TestEqual_SameResult(t *testing.T) {
	env := makeEnv("X", "1")
	r := digester.Digest(env, digester.Options{})
	if !digester.Equal(r, r) {
		t.Error("Equal should return true for identical results")
	}
}

func TestEqual_DifferentResults(t *testing.T) {
	r1 := digester.Digest(makeEnv("A", "1"), digester.Options{})
	r2 := digester.Digest(makeEnv("A", "2"), digester.Options{})
	if digester.Equal(r1, r2) {
		t.Error("Equal should return false for different fingerprints")
	}
}
