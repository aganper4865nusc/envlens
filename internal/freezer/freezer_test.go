package freezer_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/freezer"
)

func makeEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"SECRET":   "s3cr3t",
	}
}

func TestFreeze_KeysAreSorted(t *testing.T) {
	f := freezer.Freeze(makeEnv())
	keys := f.Keys()
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}

func TestFreeze_GetKnownKey(t *testing.T) {
	f := freezer.Freeze(makeEnv())
	v, ok := f.Get("DB_HOST")
	if !ok {
		t.Fatal("expected DB_HOST to be present")
	}
	if v != "localhost" {
		t.Errorf("expected localhost, got %q", v)
	}
}

func TestFreeze_GetMissingKey(t *testing.T) {
	f := freezer.Freeze(makeEnv())
	_, ok := f.Get("MISSING")
	if ok {
		t.Error("expected MISSING key to be absent")
	}
}

func TestFreeze_Len(t *testing.T) {
	f := freezer.Freeze(makeEnv())
	if f.Len() != 4 {
		t.Errorf("expected 4 keys, got %d", f.Len())
	}
}

func TestFreeze_Has(t *testing.T) {
	f := freezer.Freeze(makeEnv())
	if !f.Has("SECRET") {
		t.Error("expected Has(SECRET) to be true")
	}
	if f.Has("NOT_THERE") {
		t.Error("expected Has(NOT_THERE) to be false")
	}
}

func TestThaw_ReturnsMutableCopy(t *testing.T) {
	f := freezer.Freeze(makeEnv())
	thawed := f.Thaw()
	thawed["NEW_KEY"] = "added"

	if f.Has("NEW_KEY") {
		t.Error("mutation of thawed map should not affect frozen env")
	}
	if _, ok := thawed["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY in thawed map")
	}
}

func TestFreeze_OriginalMapUnmodified(t *testing.T) {
	env := makeEnv()
	_ = freezer.Freeze(env)
	env["INJECTED"] = "yes"

	f2 := freezer.Freeze(env)
	if !f2.Has("INJECTED") {
		t.Error("second freeze should see INJECTED key")
	}
}

func TestFreeze_EmptyEnv(t *testing.T) {
	f := freezer.Freeze(map[string]string{})
	if f.Len() != 0 {
		t.Errorf("expected 0 keys, got %d", f.Len())
	}
	if len(f.Keys()) != 0 {
		t.Error("expected empty keys slice")
	}
}
