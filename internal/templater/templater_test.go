package templater_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/templater"
)

func TestRender_BasicSubstitution(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	tmpl := "host={{ .APP_HOST }} port={{ .APP_PORT }}"
	res, err := templater.Render(tmpl, env, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "host=localhost port=8080" {
		t.Errorf("got %q", res.Output)
	}
}

func TestRender_DefaultFunc(t *testing.T) {
	env := map[string]string{"APP_HOST": ""}
	tmpl := `{{ default "127.0.0.1" .APP_HOST }}`
	res, err := templater.Render(tmpl, env, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "127.0.0.1" {
		t.Errorf("got %q", res.Output)
	}
}

func TestRender_UpperLowerFuncs(t *testing.T) {
	env := map[string]string{"ENV": "production"}
	tmpl := `{{ upper .ENV }}/{{ lower .ENV }}`
	res, err := templater.Render(tmpl, env, templater.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "PRODUCTION/production" {
		t.Errorf("got %q", res.Output)
	}
}

func TestRender_MissingKey_Error(t *testing.T) {
	env := map[string]string{}
	tmpl := "{{ .MISSING_KEY }}"
	_, err := templater.Render(tmpl, env, templater.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRender_MissingKey_Zero(t *testing.T) {
	env := map[string]string{}
	tmpl := "value={{ .MISSING_KEY }}"
	opts := templater.DefaultOptions()
	opts.MissingKey = "zero"
	res, err := templater.Render(tmpl, env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "value=" {
		t.Errorf("got %q", res.Output)
	}
}

func TestRender_WarnFunc_CollectsWarnings(t *testing.T) {
	env := map[string]string{}
	tmpl := `{{ warn "deprecated key used" }}ok`
	opts := templater.DefaultOptions()
	opts.MissingKey = "zero"
	res, err := templater.Render(tmpl, env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Warnings) != 1 || res.Warnings[0] != "deprecated key used" {
		t.Errorf("expected 1 warning, got %v", res.Warnings)
	}
	if !strings.Contains(res.Output, "ok") {
		t.Errorf("output should contain 'ok', got %q", res.Output)
	}
}

func TestRender_CustomDelimiters(t *testing.T) {
	env := map[string]string{"DB": "postgres"}
	tmpl := "db=<< .DB >>"
	opts := templater.Options{LeftDelim: "<<", RightDelim: ">>", MissingKey: "error"}
	res, err := templater.Render(tmpl, env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "db=postgres" {
		t.Errorf("got %q", res.Output)
	}
}

func TestRender_InvalidTemplate_ReturnsError(t *testing.T) {
	env := map[string]string{}
	tmpl := "{{ .UNCLOSED"
	_, err := templater.Render(tmpl, env, templater.DefaultOptions())
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}
