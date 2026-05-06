package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envlens")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_MissingTarget(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when --target missing")
	}
	if !strings.Contains(string(out), "--target is required") {
		t.Errorf("expected error message, got: %s", out)
	}
}

func TestMain_AuditOnlyTextOutput(t *testing.T) {
	bin := buildBinary(t)
	target := writeTempFile(t, "APP_KEY=plaintext_secret\nDB_HOST=localhost\n")
	cmd := exec.Command(bin, "--target", target, "--audit", "--format", "text")
	out, err := cmd.CombinedOutput()
	// audit issues do not cause non-zero exit
	if err != nil {
		t.Fatalf("unexpected exit error: %v\noutput: %s", err, out)
	}
	output := string(out)
	if !strings.Contains(output, "Audit") {
		t.Errorf("expected audit section in output, got: %s", output)
	}
}

func TestMain_DiffTextOutput(t *testing.T) {
	bin := buildBinary(t)
	base := writeTempFile(t, "FOO=bar\nBAZ=qux\n")
	target := writeTempFile(t, "FOO=bar\nNEW_KEY=hello\n")
	cmd := exec.Command(bin, "--base", base, "--target", target, "--format", "text")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected exit error: %v\noutput: %s", err, out)
	}
	output := string(out)
	if !strings.Contains(output, "Diff") {
		t.Errorf("expected diff section in output, got: %s", output)
	}
}
