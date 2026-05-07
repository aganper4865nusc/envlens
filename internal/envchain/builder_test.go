package envchain_test

import (
	"testing"

	"github.com/user/envlens/internal/envchain"
)

func TestStripPrefixStep(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "OTHER": "x"}
	c := envchain.New()
	c.Add(envchain.StripPrefixStep("APP_").Name, envchain.StripPrefixStep("APP_").Fn)

	res, err := c.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["HOST"]; !ok {
		t.Error("expected key HOST after stripping APP_ prefix")
	}
	if _, ok := res.Env["OTHER"]; !ok {
		t.Error("expected key OTHER to remain unchanged")
	}
}

func TestAddPrefixStep(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	step := envchain.AddPrefixStep("SVC_")
	c := envchain.New().Add(step.Name, step.Fn)

	res, err := c.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST=localhost, got %q", res.Env["SVC_HOST"])
	}
	if _, ok := res.Env["HOST"]; ok {
		t.Error("original key HOST should not exist after prefix added")
	}
}

func TestDropEmptyStep(t *testing.T) {
	env := map[string]string{"A": "value", "B": "", "C": "   "}
	step := envchain.DropEmptyStep()
	c := envchain.New().Add(step.Name, step.Fn)

	res, err := c.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Env["B"]; ok {
		t.Error("expected empty key B to be dropped")
	}
	if _, ok := res.Env["C"]; ok {
		t.Error("expected blank-only key C to be dropped")
	}
	if res.Env["A"] != "value" {
		t.Error("expected key A to remain")
	}
}

func TestUppercaseKeysStep(t *testing.T) {
	env := map[string]string{"host": "localhost", "Port": "9000"}
	step := envchain.UppercaseKeysStep()
	c := envchain.New().Add(step.Name, step.Fn)

	res, err := c.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", res.Env["HOST"])
	}
	if res.Env["PORT"] != "9000" {
		t.Errorf("expected PORT=9000, got %q", res.Env["PORT"])
	}
}

func TestChain_ComposedBuilderSteps(t *testing.T) {
	env := map[string]string{"app_host": "localhost", "app_empty": "", "other": "val"}

	s1 := envchain.UppercaseKeysStep()
	s2 := envchain.StripPrefixStep("APP_")
	s3 := envchain.DropEmptyStep()

	c := envchain.New().
		Add(s1.Name, s1.Fn).
		Add(s2.Name, s2.Fn).
		Add(s3.Name, s3.Fn)

	res, err := c.Run(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", res.Env["HOST"])
	}
	if _, ok := res.Env["EMPTY"]; ok {
		t.Error("expected EMPTY key to be dropped")
	}
	if len(res.Log) != 3 {
		t.Errorf("expected 3 log entries, got %d", len(res.Log))
	}
}
