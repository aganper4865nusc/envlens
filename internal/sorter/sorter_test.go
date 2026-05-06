package sorter_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/sorter"
)

func makeEnv(keys ...string) map[string]string {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = "value"
	}
	return m
}

func TestSort_Alphabetical(t *testing.T) {
	env := makeEnv("ZEBRA", "APPLE", "MANGO")
	res := sorter.Sort(env, sorter.Alphabetical)
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range res.Keys {
		if k != expected[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestSort_AlphabeticalDesc(t *testing.T) {
	env := makeEnv("ZEBRA", "APPLE", "MANGO")
	res := sorter.Sort(env, sorter.AlphabeticalDesc)
	expected := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range res.Keys {
		if k != expected[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestSort_ByPrefix_GroupsCorrectly(t *testing.T) {
	env := makeEnv("DB_HOST", "DB_PORT", "APP_NAME", "APP_ENV", "PORT")
	res := sorter.Sort(env, sorter.ByPrefix)

	if len(res.Groups) == 0 {
		t.Fatal("expected non-empty groups")
	}

	dbKeys := res.Groups["DB"]
	if len(dbKeys) != 2 {
		t.Errorf("DB group: got %d keys, want 2", len(dbKeys))
	}

	appKeys := res.Groups["APP"]
	if len(appKeys) != 2 {
		t.Errorf("APP group: got %d keys, want 2", len(appKeys))
	}

	noPrefix := res.Groups[""]
	if len(noPrefix) != 1 || noPrefix[0] != "PORT" {
		t.Errorf("no-prefix group: got %v, want [PORT]", noPrefix)
	}
}

func TestSort_ByPrefix_KeysStillSorted(t *testing.T) {
	env := makeEnv("Z_ONE", "A_TWO", "M_THREE")
	res := sorter.Sort(env, sorter.ByPrefix)
	if len(res.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(res.Keys))
	}
	for i := 1; i < len(res.Keys); i++ {
		if res.Keys[i] < res.Keys[i-1] {
			t.Errorf("keys not sorted at index %d: %q < %q", i, res.Keys[i], res.Keys[i-1])
		}
	}
}

func TestSort_EmptyMap(t *testing.T) {
	res := sorter.Sort(map[string]string{}, sorter.Alphabetical)
	if len(res.Keys) != 0 {
		t.Errorf("expected empty keys, got %v", res.Keys)
	}
}

func TestSort_AlphabeticalGroupsNil(t *testing.T) {
	env := makeEnv("A", "B")
	res := sorter.Sort(env, sorter.Alphabetical)
	if res.Groups != nil {
		t.Errorf("expected nil Groups for Alphabetical mode, got %v", res.Groups)
	}
}
