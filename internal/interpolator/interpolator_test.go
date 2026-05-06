package interpolator_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/interpolator"
)

func TestInterpolate_BraceStyle(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"URL":  "http://${HOST}:8080",
	}
	res := interpolator.Interpolate(env, nil)
	if got := res.Env["URL"]; got != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %s", got)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", res.Warnings)
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	env := map[string]string{
		"USER": "alice",
		"GREETING": "hello $USER",
	}
	res := interpolator.Interpolate(env, nil)
	if got := res.Env["GREETING"]; got != "hello alice" {
		t.Errorf("expected 'hello alice', got %s", got)
	}
}

func TestInterpolate_UndefinedVariable_Warns(t *testing.T) {
	env := map[string]string{
		"PATH_VAL": "/home/${UNDEFINED_VAR}/bin",
	}
	res := interpolator.Interpolate(env, nil)
	if got := res.Env["PATH_VAL"]; got != "/home//bin" {
		t.Errorf("expected '/home//bin', got %s", got)
	}
	if len(res.Warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(res.Warnings))
	}
	if res.Warnings[0] != "undefined variable: UNDEFINED_VAR" {
		t.Errorf("unexpected warning: %s", res.Warnings[0])
	}
}

func TestInterpolate_ExternalContext(t *testing.T) {
	env := map[string]string{
		"DSN": "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}/mydb",
	}
	ctx := map[string]string{
		"DB_USER": "root",
		"DB_PASS": "secret",
		"DB_HOST": "db.internal",
	}
	res := interpolator.Interpolate(env, ctx)
	want := "postgres://root:secret@db.internal/mydb"
	if got := res.Env["DSN"]; got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}

func TestInterpolate_NoReferences_Unchanged(t *testing.T) {
	env := map[string]string{
		"PLAIN": "no-references-here",
	}
	res := interpolator.Interpolate(env, nil)
	if got := res.Env["PLAIN"]; got != "no-references-here" {
		t.Errorf("expected unchanged value, got %s", got)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", res.Warnings)
	}
}

func TestInterpolate_OriginalMapUnmodified(t *testing.T) {
	env := map[string]string{
		"BASE": "base",
		"DERIVED": "${BASE}_suffix",
	}
	_ = interpolator.Interpolate(env, nil)
	if env["DERIVED"] != "${BASE}_suffix" {
		t.Error("original map should not be modified")
	}
}
