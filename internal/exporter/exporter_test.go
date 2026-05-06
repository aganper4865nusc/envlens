package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/exporter"
)

func TestExport_ShellFormat(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"DB_PASS": `sec"ret`,
	}
	var buf bytes.Buffer
	if err := exporter.Export(&buf, env, exporter.FormatShell); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `export APP_ENV="production"`) {
		t.Errorf("missing APP_ENV shell export, got:\n%s", out)
	}
	if !strings.Contains(out, `export DB_PASS="sec\"ret"`) {
		t.Errorf("expected escaped quote in DB_PASS, got:\n%s", out)
	}
}

func TestExport_DockerFormat(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	var buf bytes.Buffer
	if err := exporter.Export(&buf, env, exporter.FormatDocker); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("missing HOST in docker format, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=5432") {
		t.Errorf("missing PORT in docker format, got:\n%s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	var buf bytes.Buffer
	if err := exporter.Export(&buf, env, exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"KEY": "val"`) {
		t.Errorf("expected JSON output, got:\n%s", out)
	}
}

func TestExport_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	var buf bytes.Buffer
	_ = exporter.Export(&buf, env, exporter.FormatDocker)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected A_KEY first, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected Z_KEY last, got %s", lines[2])
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, map[string]string{"K": "v"}, exporter.Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestExport_EmptyEnv(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(&buf, map[string]string{}, exporter.FormatShell); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty env, got: %s", buf.String())
	}
}
