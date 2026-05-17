package mapper_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/mapper"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestMap_RenamesKey(t *testing.T) {
	env := makeEnv("OLD_KEY", "value")
	res := mapper.Map(env, mapper.Options{
		KeyRules: []mapper.KeyRule{{From: "OLD_KEY", To: "NEW_KEY"}},
	})
	if _, ok := res.Env["NEW_KEY"]; !ok {
		t.Fatal("expected NEW_KEY in result")
	}
	if _, ok := res.Env["OLD_KEY"]; ok {
		t.Fatal("expected OLD_KEY to be removed")
	}
	if len(res.RenamedKeys) != 1 || res.RenamedKeys[0] != "OLD_KEY" {
		t.Errorf("unexpected RenamedKeys: %v", res.RenamedKeys)
	}
}

func TestMap_SubstitutesValue(t *testing.T) {
	env := makeEnv("ENV", "staging")
	res := mapper.Map(env, mapper.Options{
		ValueRules: []mapper.ValueRule{{Key: "ENV", From: "staging", To: "production"}},
	})
	if res.Env["ENV"] != "production" {
		t.Errorf("expected production, got %q", res.Env["ENV"])
	}
	if len(res.Substituted) != 1 || res.Substituted[0] != "ENV" {
		t.Errorf("unexpected Substituted: %v", res.Substituted)
	}
}

func TestMap_DropUnmapped(t *testing.T) {
	env := makeEnv("KEEP", "yes", "DROP", "no")
	res := mapper.Map(env, mapper.Options{
		KeyRules:     []mapper.KeyRule{{From: "KEEP", To: "KEEP"}},
		DropUnmapped: true,
	})
	if _, ok := res.Env["DROP"]; ok {
		t.Fatal("expected DROP to be removed")
	}
	if _, ok := res.Env["KEEP"]; !ok {
		t.Fatal("expected KEEP to be present")
	}
}

func TestMap_NoRules_PassThrough(t *testing.T) {
	env := makeEnv("A", "1", "B", "2")
	res := mapper.Map(env, mapper.Options{})
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
	if len(res.RenamedKeys) != 0 || len(res.Substituted) != 0 {
		t.Error("expected no renames or substitutions")
	}
}

func TestMap_RenameAndSubstituteComposed(t *testing.T) {
	env := makeEnv("LOG_LEVEL", "debug")
	res := mapper.Map(env, mapper.Options{
		KeyRules:   []mapper.KeyRule{{From: "LOG_LEVEL", To: "LOGGING_LEVEL"}},
		ValueRules: []mapper.ValueRule{{Key: "LOG_LEVEL", From: "debug", To: "DEBUG"}},
	})
	if res.Env["LOGGING_LEVEL"] != "DEBUG" {
		t.Errorf("expected DEBUG, got %q", res.Env["LOGGING_LEVEL"])
	}
}

func TestMap_OriginalMapUnmodified(t *testing.T) {
	env := makeEnv("K", "v")
	mapper.Map(env, mapper.Options{
		KeyRules: []mapper.KeyRule{{From: "K", To: "K2"}},
	})
	if _, ok := env["K"]; !ok {
		t.Error("original map should not be modified")
	}
}
