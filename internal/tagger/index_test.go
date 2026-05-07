package tagger_test

import (
	"reflect"
	"testing"

	"github.com/yourorg/envlens/internal/tagger"
)

func TestIndex_GroupsByTag(t *testing.T) {
	env := makeEnv("DB_HOST", "h", "DB_PORT", "5432", "APP_ENV", "prod")
	rules := []tagger.Rule{
		{Tag: "database", Prefixes: []string{"DB_"}},
		{Tag: "app", Prefixes: []string{"APP_"}},
	}
	results := tagger.Tag(env, rules)
	idx := tagger.Index(results)

	if !reflect.DeepEqual(idx["database"], []string{"DB_HOST", "DB_PORT"}) {
		t.Errorf("unexpected database keys: %v", idx["database"])
	}
	if !reflect.DeepEqual(idx["app"], []string{"APP_ENV"}) {
		t.Errorf("unexpected app keys: %v", idx["app"])
	}
}

func TestIndex_EmptyResults(t *testing.T) {
	idx := tagger.Index([]tagger.Result{})
	if len(idx) != 0 {
		t.Errorf("expected empty index, got %v", idx)
	}
}

func TestIndex_KeysAreSorted(t *testing.T) {
	env := makeEnv("Z_KEY", "1", "A_KEY", "2", "M_KEY", "3")
	rules := []tagger.Rule{{Tag: "all", Pattern: `_KEY$`}}
	results := tagger.Tag(env, rules)
	idx := tagger.Index(results)
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	if !reflect.DeepEqual(idx["all"], expected) {
		t.Errorf("expected sorted keys %v, got %v", expected, idx["all"])
	}
}

func TestTagsFor_Found(t *testing.T) {
	results := []tagger.Result{
		{Key: "DB_HOST", Tags: []string{"database", "infra"}},
		{Key: "APP_ENV", Tags: []string{"app"}},
	}
	tags := tagger.TagsFor("DB_HOST", results)
	if !reflect.DeepEqual(tags, []string{"database", "infra"}) {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestTagsFor_NotFound(t *testing.T) {
	results := []tagger.Result{
		{Key: "DB_HOST", Tags: []string{"database"}},
	}
	tags := tagger.TagsFor("MISSING", results)
	if tags != nil {
		t.Errorf("expected nil for missing key, got %v", tags)
	}
}

func TestTagsFor_NoTagsKey(t *testing.T) {
	results := []tagger.Result{
		{Key: "LOG_LEVEL", Tags: []string{}},
	}
	tags := tagger.TagsFor("LOG_LEVEL", results)
	if len(tags) != 0 {
		t.Errorf("expected empty tag slice, got %v", tags)
	}
}
