package normalizer_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/normalizer"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	env := map[string]string{"app_name": "envlens", "Port": "8080"}
	res := normalizer.Normalize(env, normalizer.Options{UppercaseKeys: true})
	if _, ok := res.Env["APP_NAME"]; !ok {
		t.Error("expected APP_NAME key")
	}
	if _, ok := res.Env["PORT"]; !ok {
		t.Error("expected PORT key")
	}
	if len(res.RenamedKeys) != 2 {
		t.Errorf("expected 2 renamed keys, got %d", len(res.RenamedKeys))
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	env := map[string]string{"KEY": "  hello  ", "OTHER": "clean"}
	res := normalizer.Normalize(env, normalizer.Options{TrimValues: true})
	if res.Env["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", res.Env["KEY"])
	}
	if len(res.TrimmedValues) != 1 || res.TrimmedValues[0] != "KEY" {
		t.Errorf("expected TrimmedValues=[KEY], got %v", res.TrimmedValues)
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	env := map[string]string{"PRESENT": "value", "EMPTY": ""}
	res := normalizer.Normalize(env, normalizer.Options{RemoveEmpty: true})
	if _, ok := res.Env["EMPTY"]; ok {
		t.Error("expected EMPTY to be dropped")
	}
	if len(res.DroppedKeys) != 1 || res.DroppedKeys[0] != "EMPTY" {
		t.Errorf("expected DroppedKeys=[EMPTY], got %v", res.DroppedKeys)
	}
	if _, ok := res.Env["PRESENT"]; !ok {
		t.Error("expected PRESENT to remain")
	}
}

func TestNormalize_TrimKeys(t *testing.T) {
	env := map[string]string{"  SPACED  ": "val"}
	res := normalizer.Normalize(env, normalizer.Options{TrimKeys: true})
	if _, ok := res.Env["SPACED"]; !ok {
		t.Error("expected trimmed key 'SPACED'")
	}
}

func TestNormalize_NoChanges(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	res := normalizer.Normalize(env, normalizer.Options{})
	if res.Env["KEY"] != "value" {
		t.Errorf("expected value unchanged, got %q", res.Env["KEY"])
	}
	if len(res.RenamedKeys) != 0 || len(res.TrimmedValues) != 0 || len(res.DroppedKeys) != 0 {
		t.Error("expected no changes")
	}
}

func TestNormalize_OriginalMapUnmodified(t *testing.T) {
	env := map[string]string{"lower": "  spaced  "}
	normalizer.Normalize(env, normalizer.DefaultOptions())
	if env["lower"] != "  spaced  " {
		t.Error("original map was modified")
	}
}

func TestNormalize_DefaultOptions(t *testing.T) {
	env := map[string]string{"my_key": " trimmed "}
	res := normalizer.Normalize(env, normalizer.DefaultOptions())
	if _, ok := res.Env["MY_KEY"]; !ok {
		t.Error("expected MY_KEY after default normalization")
	}
	if res.Env["MY_KEY"] != "trimmed" {
		t.Errorf("expected value 'trimmed', got %q", res.Env["MY_KEY"])
	}
}
