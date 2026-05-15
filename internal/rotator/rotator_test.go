package rotator_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/rotator"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestRotate_BasicRename(t *testing.T) {
	env := makeEnv("OLD_KEY", "hello")
	ops := []rotator.Op{{FromKey: "OLD_KEY", ToKey: "NEW_KEY"}}
	res, err := rotator.Rotate(env, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if res.Env["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", res.Env["NEW_KEY"])
	}
	if len(res.Rotated) != 1 || res.Rotated[0] != "OLD_KEY" {
		t.Errorf("expected Rotated=[OLD_KEY], got %v", res.Rotated)
	}
}

func TestRotate_MissingFromKey_Skipped(t *testing.T) {
	env := makeEnv("PRESENT", "val")
	ops := []rotator.Op{{FromKey: "MISSING", ToKey: "NEW"}}
	res, err := rotator.Rotate(env, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected Skipped=[MISSING], got %v", res.Skipped)
	}
}

func TestRotate_ConflictRecorded(t *testing.T) {
	env := makeEnv("A", "1", "B", "2")
	ops := []rotator.Op{{FromKey: "A", ToKey: "B"}}
	res, err := rotator.Rotate(env, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflict) != 1 || res.Conflict[0] != "B" {
		t.Errorf("expected Conflict=[B], got %v", res.Conflict)
	}
	if res.Env["B"] != "1" {
		t.Errorf("expected B=1 (overwritten), got %q", res.Env["B"])
	}
}

func TestRotate_WithTransform(t *testing.T) {
	env := makeEnv("SECRET", "mysecret")
	ops := []rotator.Op{
		{FromKey: "SECRET", ToKey: "APP_SECRET", Transform: strings.ToUpper},
	}
	res, err := rotator.Rotate(env, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["APP_SECRET"] != "MYSECRET" {
		t.Errorf("expected APP_SECRET=MYSECRET, got %q", res.Env["APP_SECRET"])
	}
}

func TestRotate_OriginalMapUnmodified(t *testing.T) {
	env := makeEnv("OLD", "value")
	ops := []rotator.Op{{FromKey: "OLD", ToKey: "NEW"}}
	_, err := rotator.Rotate(env, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["OLD"]; !ok {
		t.Error("original map should not be modified")
	}
}

func TestRotate_EmptyKeyError(t *testing.T) {
	env := makeEnv("A", "1")
	ops := []rotator.Op{{FromKey: "", ToKey: "B"}}
	_, err := rotator.Rotate(env, ops)
	if err == nil {
		t.Error("expected error for empty FromKey")
	}
}

func TestRotate_WhitespaceInToKeyError(t *testing.T) {
	env := makeEnv("A", "1")
	ops := []rotator.Op{{FromKey: "A", ToKey: "B C"}}
	_, err := rotator.Rotate(env, ops)
	if err == nil {
		t.Error("expected error for whitespace in ToKey")
	}
}
