package patcher_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/patcher"
)

func baseEnv() map[string]string {
	return map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
		"DEBUG": "true",
	}
}

func TestPatch_SetNewKey(t *testing.T) {
	out, res := patcher.Patch(baseEnv(), []patcher.Op{
		{Type: patcher.OpSet, Key: "TIMEOUT", Value: "30s"},
	})
	if out["TIMEOUT"] != "30s" {
		t.Errorf("expected TIMEOUT=30s, got %s", out["TIMEOUT"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestPatch_SetOverwritesExisting(t *testing.T) {
	out, _ := patcher.Patch(baseEnv(), []patcher.Op{
		{Type: patcher.OpSet, Key: "PORT", Value: "9090"},
	})
	if out["PORT"] != "9090" {
		t.Errorf("expected PORT=9090, got %s", out["PORT"])
	}
}

func TestPatch_UnsetExistingKey(t *testing.T) {
	out, res := patcher.Patch(baseEnv(), []patcher.Op{
		{Type: patcher.OpUnset, Key: "DEBUG"},
	})
	if _, ok := out["DEBUG"]; ok {
		t.Error("expected DEBUG to be removed")
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestPatch_UnsetMissingKey_Skipped(t *testing.T) {
	_, res := patcher.Patch(baseEnv(), []patcher.Op{
		{Type: patcher.OpUnset, Key: "MISSING"},
	})
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestPatch_RenameKey(t *testing.T) {
	out, res := patcher.Patch(baseEnv(), []patcher.Op{
		{Type: patcher.OpRename, Key: "HOST", To: "DB_HOST"},
	})
	if _, ok := out["HOST"]; ok {
		t.Error("old key HOST should not exist")
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", out["DB_HOST"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestPatch_RenameToExistingKey_Warns(t *testing.T) {
	_, res := patcher.Patch(baseEnv(), []patcher.Op{
		{Type: patcher.OpRename, Key: "HOST", To: "PORT"},
	})
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(res.Warnings))
	}
}

func TestPatch_OriginalMapUnmodified(t *testing.T) {
	env := baseEnv()
	patcher.Patch(env, []patcher.Op{
		{Type: patcher.OpSet, Key: "NEW", Value: "val"},
		{Type: patcher.OpUnset, Key: "PORT"},
	})
	if _, ok := env["NEW"]; ok {
		t.Error("original map should not have NEW")
	}
	if env["PORT"] != "8080" {
		t.Error("original map PORT should be unchanged")
	}
}

func TestPatch_UnknownOpType_Warns(t *testing.T) {
	_, res := patcher.Patch(baseEnv(), []patcher.Op{
		{Type: "invalid", Key: "X"},
	})
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning for unknown op, got %d", len(res.Warnings))
	}
}
