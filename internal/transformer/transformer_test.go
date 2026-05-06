package transformer_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/transformer"
)

func TestTransform_AddPrefix(t *testing.T) {
	src := map[string]string{"HOST": "localhost", "PORT": "8080"}
	res := transformer.Transform(src, transformer.Options{AddPrefix: "APP_"})

	if _, ok := res.Env["APP_HOST"]; !ok {
		t.Error("expected APP_HOST key")
	}
	if _, ok := res.Env["APP_PORT"]; !ok {
		t.Error("expected APP_PORT key")
	}
	if len(res.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(res.Changes))
	}
}

func TestTransform_StripPrefix(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "OTHER": "x"}
	res := transformer.Transform(src, transformer.Options{StripPrefix: "APP_"})

	if _, ok := res.Env["HOST"]; !ok {
		t.Error("expected HOST after stripping prefix")
	}
	if _, ok := res.Env["OTHER"]; !ok {
		t.Error("expected OTHER unchanged")
	}
}

func TestTransform_UppercaseKeys(t *testing.T) {
	src := map[string]string{"db_host": "localhost"}
	res := transformer.Transform(src, transformer.Options{UppercaseKeys: true})

	if v, ok := res.Env["DB_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got Env=%v", res.Env)
	}
}

func TestTransform_UppercaseValues(t *testing.T) {
	src := map[string]string{"ENV": "production"}
	res := transformer.Transform(src, transformer.Options{UppercaseValues: true})

	if res.Env["ENV"] != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %s", res.Env["ENV"])
	}
}

func TestTransform_NoChanges(t *testing.T) {
	src := map[string]string{"KEY": "value"}
	res := transformer.Transform(src, transformer.Options{})

	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(res.Changes))
	}
	if res.Env["KEY"] != "value" {
		t.Error("expected KEY=value unchanged")
	}
}

func TestTransform_OriginalUnmodified(t *testing.T) {
	src := map[string]string{"host": "localhost"}
	transformer.Transform(src, transformer.Options{AddPrefix: "X_", UppercaseKeys: true})

	if _, ok := src["host"]; !ok {
		t.Error("original map was modified")
	}
}

func TestTransform_ChangeRecordsOldKey(t *testing.T) {
	src := map[string]string{"db_pass": "secret"}
	res := transformer.Transform(src, transformer.Options{UppercaseKeys: true})

	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	c := res.Changes[0]
	if c.OldKey != "db_pass" {
		t.Errorf("expected OldKey=db_pass, got %s", c.OldKey)
	}
	if c.Key != "DB_PASS" {
		t.Errorf("expected Key=DB_PASS, got %s", c.Key)
	}
}

func TestTransform_AddPrefixAndUppercaseKeys(t *testing.T) {
	// Verify that combining AddPrefix and UppercaseKeys applies both transformations.
	src := map[string]string{"db_host": "localhost"}
	res := transformer.Transform(src, transformer.Options{AddPrefix: "APP_", UppercaseKeys: true})

	if _, ok := res.Env["APP_DB_HOST"]; !ok {
		t.Errorf("expected APP_DB_HOST key, got Env=%v", res.Env)
	}
	if res.Env["APP_DB_HOST"] != "localhost" {
		t.Errorf("expected value localhost, got %s", res.Env["APP_DB_HOST"])
	}
}
