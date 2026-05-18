package pipeline_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/pipeline"
)

func TestUppercaseKeysStage(t *testing.T) {
	p := pipeline.New().Add(pipeline.UppercaseKeysStage())
	results, err := p.Run(makeEnv("foo", "bar", "baz", "qux"))
	if err != nil {
		t.Fatal(err)
	}
	final := pipeline.Final(results)
	if _, ok := final["FOO"]; !ok {
		t.Error("expected FOO key")
	}
	if _, ok := final["BAZ"]; !ok {
		t.Error("expected BAZ key")
	}
}

func TestAddPrefixStage(t *testing.T) {
	p := pipeline.New().Add(pipeline.AddPrefixStage("APP_"))
	results, err := p.Run(makeEnv("NAME", "envlens"))
	if err != nil {
		t.Fatal(err)
	}
	final := pipeline.Final(results)
	if v, ok := final["APP_NAME"]; !ok || v != "envlens" {
		t.Errorf("expected APP_NAME=envlens, got %v", final)
	}
}

func TestStripPrefixStage(t *testing.T) {
	p := pipeline.New().Add(pipeline.StripPrefixStage("APP_"))
	results, err := p.Run(makeEnv("APP_NAME", "envlens", "APP_ENV", "prod"))
	if err != nil {
		t.Fatal(err)
	}
	final := pipeline.Final(results)
	if _, ok := final["NAME"]; !ok {
		t.Error("expected NAME after strip")
	}
}

func TestDropEmptyStage(t *testing.T) {
	p := pipeline.New().Add(pipeline.DropEmptyStage())
	results, err := p.Run(makeEnv("KEY", "value", "EMPTY", "", "BLANK", "   "))
	if err != nil {
		t.Fatal(err)
	}
	final := pipeline.Final(results)
	if _, ok := final["EMPTY"]; ok {
		t.Error("EMPTY should have been dropped")
	}
	if _, ok := final["BLANK"]; ok {
		t.Error("BLANK should have been dropped")
	}
	if _, ok := final["KEY"]; !ok {
		t.Error("KEY should be present")
	}
}

func TestComposedBuilderStages(t *testing.T) {
	p := pipeline.New().
		Add(pipeline.DropEmptyStage()).
		Add(pipeline.UppercaseKeysStage()).
		Add(pipeline.AddPrefixStage("CI_"))
	results, err := p.Run(makeEnv("name", "envlens", "empty", ""))
	if err != nil {
		t.Fatal(err)
	}
	final := pipeline.Final(results)
	if v, ok := final["CI_NAME"]; !ok || v != "envlens" {
		t.Errorf("expected CI_NAME=envlens, got %v", final)
	}
	if _, ok := final["CI_EMPTY"]; ok {
		t.Error("CI_EMPTY should have been dropped")
	}
}
