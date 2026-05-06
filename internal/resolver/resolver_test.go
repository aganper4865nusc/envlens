package resolver_test

import (
	"os"
	"testing"

	"github.com/yourorg/envlens/internal/resolver"
)

func TestResolve_FileSource(t *testing.T) {
	fileEnv := map[string]string{"APP_PORT": "8080"}
	results, _ := resolver.Resolve(fileEnv, resolver.Options{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Source != "file" {
		t.Errorf("expected source=file, got %s", results[0].Source)
	}
	if results[0].Value != "8080" {
		t.Errorf("expected value=8080, got %s", results[0].Value)
	}
}

func TestResolve_EnvOverride(t *testing.T) {
	t.Setenv("APP_PORT", "9090")
	fileEnv := map[string]string{"APP_PORT": "8080"}
	results, _ := resolver.Resolve(fileEnv, resolver.Options{AllowEnvOverride: true})
	if results[0].Value != "9090" {
		t.Errorf("expected overridden value 9090, got %s", results[0].Value)
	}
	if !results[0].Override {
		t.Error("expected Override=true")
	}
	if results[0].Source != "env" {
		t.Errorf("expected source=env, got %s", results[0].Source)
	}
}

func TestResolve_NoOverrideWhenDisabled(t *testing.T) {
	t.Setenv("APP_PORT", "9090")
	fileEnv := map[string]string{"APP_PORT": "8080"}
	results, _ := resolver.Resolve(fileEnv, resolver.Options{AllowEnvOverride: false})
	if results[0].Value != "8080" {
		t.Errorf("expected file value 8080, got %s", results[0].Value)
	}
}

func TestResolve_DefaultFallback(t *testing.T) {
	fileEnv := map[string]string{}
	opts := resolver.Options{
		Defaults: map[string]string{"LOG_LEVEL": "info"},
	}
	results, _ := resolver.Resolve(fileEnv, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Source != "default" {
		t.Errorf("expected source=default, got %s", results[0].Source)
	}
	if results[0].Value != "info" {
		t.Errorf("expected value=info, got %s", results[0].Value)
	}
}

func TestResolve_FileWinsOverDefault(t *testing.T) {
	fileEnv := map[string]string{"LOG_LEVEL": "debug"}
	opts := resolver.Options{
		Defaults: map[string]string{"LOG_LEVEL": "info"},
	}
	results, _ := resolver.Resolve(fileEnv, opts)
	if results[0].Value != "debug" {
		t.Errorf("expected file value debug, got %s", results[0].Value)
	}
	if results[0].Source != "file" {
		t.Errorf("expected source=file, got %s", results[0].Source)
	}
}

func TestToMap(t *testing.T) {
	resolutions := []resolver.Resolution{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	m := resolver.ToMap(resolutions)
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestResolve_EmptyDefaultFlagged(t *testing.T) {
	os.Unsetenv("EMPTY_KEY")
	fileEnv := map[string]string{}
	opts := resolver.Options{
		Defaults: map[string]string{"EMPTY_KEY": ""},
	}
	_, missing := resolver.Resolve(fileEnv, opts)
	if len(missing) == 0 {
		t.Error("expected empty default to appear in missing list")
	}
}
