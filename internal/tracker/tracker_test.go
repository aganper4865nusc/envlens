package tracker_test

import (
	"testing"

	"github.com/user/envlens/internal/tracker"
)

func TestTrack_NoChanges(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := tracker.Track([]map[string]string{a, b})
	if r.TotalChanges != 0 {
		t.Fatalf("expected 0 changes, got %d", r.TotalChanges)
	}
	if len(r.Histories) != 0 {
		t.Fatalf("expected no histories, got %d", len(r.Histories))
	}
}

func TestTrack_SingleChange(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	r := tracker.Track([]map[string]string{a, b})
	if r.TotalChanges != 1 {
		t.Fatalf("expected 1 change, got %d", r.TotalChanges)
	}
	if len(r.Histories) != 1 {
		t.Fatalf("expected 1 history entry")
	}
	ev := r.Histories[0].Events[0]
	if ev.OldValue != "old" || ev.NewValue != "new" {
		t.Errorf("unexpected event values: %+v", ev)
	}
}

func TestTrack_AddedKey(t *testing.T) {
	a := map[string]string{}
	b := map[string]string{"NEW_KEY": "value"}
	r := tracker.Track([]map[string]string{a, b})
	if r.TotalChanges != 1 {
		t.Fatalf("expected 1 change, got %d", r.TotalChanges)
	}
	ev := r.Histories[0].Events[0]
	if ev.OldValue != "" || ev.NewValue != "value" {
		t.Errorf("unexpected event: %+v", ev)
	}
}

func TestTrack_RemovedKey(t *testing.T) {
	a := map[string]string{"GONE": "here"}
	b := map[string]string{}
	r := tracker.Track([]map[string]string{a, b})
	if r.TotalChanges != 1 {
		t.Fatalf("expected 1 change, got %d", r.TotalChanges)
	}
	ev := r.Histories[0].Events[0]
	if ev.OldValue != "here" || ev.NewValue != "" {
		t.Errorf("unexpected event: %+v", ev)
	}
}

func TestTrack_MultipleSnapshots(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	c := map[string]string{"X": "3"}
	r := tracker.Track([]map[string]string{a, b, c})
	if r.TotalChanges != 2 {
		t.Fatalf("expected 2 changes, got %d", r.TotalChanges)
	}
	if len(r.Histories[0].Events) != 2 {
		t.Fatalf("expected 2 events for key X")
	}
}

func TestTrack_KeysChangedSorted(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "1", "M": "1"}
	b := map[string]string{"Z": "2", "A": "2", "M": "2"}
	r := tracker.Track([]map[string]string{a, b})
	keys := tracker.KeysChanged(r)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestTrack_SingleSnapshot_ReturnsEmpty(t *testing.T) {
	r := tracker.Track([]map[string]string{{"FOO": "bar"}})
	if r.TotalChanges != 0 {
		t.Fatalf("expected 0 changes for single snapshot")
	}
}
