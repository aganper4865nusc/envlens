package patcher_test

import (
	"os"
	"testing"

	"github.com/yourusername/envlens/internal/patcher"
)

func writePatchFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "patch-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadOps_SetOp(t *testing.T) {
	path := writePatchFile(t, "set PORT=9090\n")
	ops, err := patcher.LoadOps(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 || ops[0].Type != patcher.OpSet || ops[0].Key != "PORT" || ops[0].Value != "9090" {
		t.Errorf("unexpected op: %+v", ops)
	}
}

func TestLoadOps_UnsetOp(t *testing.T) {
	path := writePatchFile(t, "unset DEBUG\n")
	ops, err := patcher.LoadOps(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 || ops[0].Type != patcher.OpUnset || ops[0].Key != "DEBUG" {
		t.Errorf("unexpected op: %+v", ops)
	}
}

func TestLoadOps_RenameOp(t *testing.T) {
	path := writePatchFile(t, "rename HOST DB_HOST\n")
	ops, err := patcher.LoadOps(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 || ops[0].Type != patcher.OpRename || ops[0].Key != "HOST" || ops[0].To != "DB_HOST" {
		t.Errorf("unexpected op: %+v", ops)
	}
}

func TestLoadOps_SkipsCommentsAndBlanks(t *testing.T) {
	path := writePatchFile(t, "# this is a comment\n\nset X=1\n")
	ops, err := patcher.LoadOps(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 {
		t.Errorf("expected 1 op, got %d", len(ops))
	}
}

func TestLoadOps_InvalidSetFormat(t *testing.T) {
	path := writePatchFile(t, "set NOEQUALS\n")
	_, err := patcher.LoadOps(path)
	if err == nil {
		t.Error("expected error for missing = in set")
	}
}

func TestLoadOps_UnknownDirective(t *testing.T) {
	path := writePatchFile(t, "apply SOMETHING\n")
	_, err := patcher.LoadOps(path)
	if err == nil {
		t.Error("expected error for unknown directive")
	}
}

func TestLoadOps_MissingFile(t *testing.T) {
	_, err := patcher.LoadOps("/nonexistent/patch.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
