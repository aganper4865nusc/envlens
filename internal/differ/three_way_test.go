package differ

import (
	"testing"
)

func TestThreeWay_NoConflict_BothSidesAgree(t *testing.T) {
	base := map[string]string{"A": "1"}
	left := map[string]string{"A": "2"}
	right := map[string]string{"A": "2"}

	results := ThreeWay(base, left, right, ThreeWayOptions{})
	if len(results) != 1 || results[0].Conflict {
		t.Fatal("expected no conflict when both sides agree")
	}
	if results[0].Resolution != "2" {
		t.Errorf("expected resolution=2, got %q", results[0].Resolution)
	}
}

func TestThreeWay_OnlyRightChanged(t *testing.T) {
	base := map[string]string{"X": "old"}
	left := map[string]string{"X": "old"}
	right := map[string]string{"X": "new"}

	results := ThreeWay(base, left, right, ThreeWayOptions{})
	if results[0].Conflict {
		t.Fatal("expected no conflict")
	}
	if results[0].Resolution != "new" {
		t.Errorf("expected right value, got %q", results[0].Resolution)
	}
}

func TestThreeWay_OnlyLeftChanged(t *testing.T) {
	base := map[string]string{"X": "old"}
	left := map[string]string{"X": "mine"}
	right := map[string]string{"X": "old"}

	results := ThreeWay(base, left, right, ThreeWayOptions{})
	if results[0].Conflict {
		t.Fatal("expected no conflict")
	}
	if results[0].Resolution != "mine" {
		t.Errorf("expected left value, got %q", results[0].Resolution)
	}
}

func TestThreeWay_Conflict_NoPreference(t *testing.T) {
	base := map[string]string{"K": "base"}
	left := map[string]string{"K": "left"}
	right := map[string]string{"K": "right"}

	results := ThreeWay(base, left, right, ThreeWayOptions{})
	if !results[0].Conflict {
		t.Fatal("expected conflict")
	}
	if results[0].Resolution != "" {
		t.Errorf("expected empty resolution, got %q", results[0].Resolution)
	}
}

func TestThreeWay_Conflict_PreferLeft(t *testing.T) {
	base := map[string]string{"K": "base"}
	left := map[string]string{"K": "ours"}
	right := map[string]string{"K": "theirs"}

	results := ThreeWay(base, left, right, ThreeWayOptions{PreferLeft: true})
	if !results[0].Conflict {
		t.Fatal("expected conflict to be flagged")
	}
	if results[0].Resolution != "ours" {
		t.Errorf("expected 'ours', got %q", results[0].Resolution)
	}
}

func TestThreeWay_Conflict_PreferRight(t *testing.T) {
	base := map[string]string{"K": "base"}
	left := map[string]string{"K": "ours"}
	right := map[string]string{"K": "theirs"}

	results := ThreeWay(base, left, right, ThreeWayOptions{PreferRight: true})
	if results[0].Resolution != "theirs" {
		t.Errorf("expected 'theirs', got %q", results[0].Resolution)
	}
}

func TestThreeWay_SortedOutput(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "1", "M": "1"}
	left := map[string]string{"Z": "1", "A": "1", "M": "1"}
	right := map[string]string{"Z": "1", "A": "1", "M": "1"}

	results := ThreeWay(base, left, right, ThreeWayOptions{})
	if results[0].Key != "A" || results[1].Key != "M" || results[2].Key != "Z" {
		t.Error("expected sorted output")
	}
}

func TestResolved_OmitsUnresolved(t *testing.T) {
	base := map[string]string{"A": "x", "B": "y"}
	left := map[string]string{"A": "left", "B": "y"}
	right := map[string]string{"A": "right", "B": "y"}

	results := ThreeWay(base, left, right, ThreeWayOptions{})
	resolved := Resolved(results)
	if _, ok := resolved["A"]; ok {
		t.Error("unresolved conflict should be omitted")
	}
	if resolved["B"] != "y" {
		t.Errorf("expected B=y, got %q", resolved["B"])
	}
}

func TestConflicts_ReturnsOnlyConflicts(t *testing.T) {
	base := map[string]string{"A": "x", "B": "y"}
	left := map[string]string{"A": "l", "B": "y"}
	right := map[string]string{"A": "r", "B": "y"}

	results := ThreeWay(base, left, right, ThreeWayOptions{})
	conflicts := Conflicts(results)
	if len(conflicts) != 1 || conflicts[0].Key != "A" {
		t.Errorf("expected 1 conflict on key A, got %+v", conflicts)
	}
}
