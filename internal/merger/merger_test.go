package merger

import (
	"testing"
)

func TestMerge_SingleSource(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res := Merge(src)

	if res.Merged["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", res.Merged["FOO"])
	}
	if len(res.Overrides) != 0 {
		t.Errorf("expected no overrides for single source, got %v", res.Overrides)
	}
}

func TestMerge_LaterSourceOverrides(t *testing.T) {
	base := map[string]string{"FOO": "base", "KEEP": "yes"}
	override := map[string]string{"FOO": "override"}
	res := Merge(base, override)

	if res.Merged["FOO"] != "override" {
		t.Errorf("expected FOO=override, got %s", res.Merged["FOO"])
	}
	if res.Merged["KEEP"] != "yes" {
		t.Errorf("expected KEEP=yes, got %s", res.Merged["KEEP"])
	}
}

func TestMerge_OverridesTracked(t *testing.T) {
	a := map[string]string{"FOO": "a", "BAR": "a"}
	b := map[string]string{"FOO": "b"}
	c := map[string]string{"FOO": "c", "BAR": "c"}
	res := Merge(a, b, c)

	if _, ok := res.Overrides["FOO"]; !ok {
		t.Error("expected FOO to be in overrides")
	}
	if _, ok := res.Overrides["BAR"]; !ok {
		t.Error("expected BAR to be in overrides")
	}
}

func TestMerge_NoOverrideForUniqueKeys(t *testing.T) {
	a := map[string]string{"ALPHA": "1"}
	b := map[string]string{"BETA": "2"}
	res := Merge(a, b)

	if _, ok := res.Overrides["ALPHA"]; ok {
		t.Error("ALPHA should not be in overrides")
	}
	if _, ok := res.Overrides["BETA"]; ok {
		t.Error("BETA should not be in overrides")
	}
}

func TestMerge_EmptySources(t *testing.T) {
	res := Merge()
	if len(res.Merged) != 0 {
		t.Errorf("expected empty merged map, got %v", res.Merged)
	}
	if len(res.Overrides) != 0 {
		t.Errorf("expected empty overrides, got %v", res.Overrides)
	}
}
