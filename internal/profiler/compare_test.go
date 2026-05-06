package profiler_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/profiler"
)

func makeProfile(source string, env map[string]string) profiler.Profile {
	return profiler.Analyze(source, env)
}

func TestCompare_KeyDelta(t *testing.T) {
	a := makeProfile("dev", map[string]string{"A": "1", "B": "2"})
	b := makeProfile("prod", map[string]string{"A": "1", "B": "2", "C": "3"})
	c := profiler.Compare(a, b)
	if c.KeyDelta != 1 {
		t.Errorf("expected KeyDelta=1, got %d", c.KeyDelta)
	}
}

func TestCompare_NewSensitiveKey(t *testing.T) {
	a := makeProfile("dev", map[string]string{"APP": "x"})
	b := makeProfile("prod", map[string]string{"APP": "x", "DB_PASSWORD": "secret"})
	c := profiler.Compare(a, b)
	if len(c.NewSensitive) != 1 || c.NewSensitive[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD in NewSensitive, got %v", c.NewSensitive)
	}
}

func TestCompare_LostSensitiveKey(t *testing.T) {
	a := makeProfile("dev", map[string]string{"API_KEY": "abc"})
	b := makeProfile("prod", map[string]string{"OTHER": "val"})
	c := profiler.Compare(a, b)
	if len(c.LostSensitive) != 1 || c.LostSensitive[0] != "API_KEY" {
		t.Errorf("expected API_KEY in LostSensitive, got %v", c.LostSensitive)
	}
}

func TestCompare_Notes_FewerKeys(t *testing.T) {
	a := makeProfile("dev", map[string]string{"A": "1", "B": "2", "C": "3"})
	b := makeProfile("prod", map[string]string{"A": "1"})
	c := profiler.Compare(a, b)
	if len(c.Notes) == 0 {
		t.Error("expected a note about fewer keys")
	}
}

func TestCompare_NoChanges(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	a := makeProfile("dev", env)
	b := makeProfile("prod", env)
	c := profiler.Compare(a, b)
	if c.KeyDelta != 0 || len(c.NewSensitive) != 0 || len(c.LostSensitive) != 0 {
		t.Errorf("expected no changes, got %+v", c)
	}
}
