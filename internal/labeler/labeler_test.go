package labeler_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/labeler"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestLabel_PrefixMatch(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "envlens")
	rules := []labeler.Rule{{Label: "database", Prefix: "DB_"}}
	res, err := labeler.Label(env, rules)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"DB_HOST", "DB_PORT"} {
		if len(res.Labels[key]) != 1 || res.Labels[key][0] != "database" {
			t.Errorf("expected key %q to have label 'database', got %v", key, res.Labels[key])
		}
	}
	if len(res.Labels["APP_NAME"]) != 0 {
		t.Errorf("APP_NAME should be unlabeled")
	}
}

func TestLabel_PatternMatch(t *testing.T) {
	env := makeEnv("SECRET_KEY", "abc", "API_SECRET", "xyz", "HOST", "localhost")
	rules := []labeler.Rule{{Label: "secret", Pattern: "(?i)secret"}}
	res, err := labeler.Label(env, rules)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"SECRET_KEY", "API_SECRET"} {
		if len(res.Labels[key]) == 0 || res.Labels[key][0] != "secret" {
			t.Errorf("expected %q to have label 'secret'", key)
		}
	}
	if len(res.Labels["HOST"]) != 0 {
		t.Errorf("HOST should be unlabeled")
	}
}

func TestLabel_ExplicitKeyMatch(t *testing.T) {
	env := makeEnv("REGION", "us-east-1", "ZONE", "a", "ENV", "prod")
	rules := []labeler.Rule{{Label: "infra", Keys: []string{"REGION", "ZONE"}}}
	res, err := labeler.Label(env, rules)
	if err != nil {
		t.Fatal(err)
	}
	if res.Labels["REGION"][0] != "infra" || res.Labels["ZONE"][0] != "infra" {
		t.Error("expected REGION and ZONE to be labeled 'infra'")
	}
}

func TestLabel_MultipleLabelsPerKey(t *testing.T) {
	env := makeEnv("DB_SECRET", "pass")
	rules := []labeler.Rule{
		{Label: "database", Prefix: "DB_"},
		{Label: "secret", Pattern: "SECRET"},
	}
	res, err := labeler.Label(env, rules)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Labels["DB_SECRET"]) != 2 {
		t.Errorf("expected 2 labels, got %v", res.Labels["DB_SECRET"])
	}
}

func TestLabel_Unlabeled(t *testing.T) {
	env := makeEnv("FOO", "bar", "BAZ", "qux")
	res, err := labeler.Label(env, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Unlabeled) != 2 {
		t.Errorf("expected 2 unlabeled keys, got %d", len(res.Unlabeled))
	}
}

func TestLabel_InvalidPattern_ReturnsError(t *testing.T) {
	env := makeEnv("KEY", "val")
	rules := []labeler.Rule{{Label: "bad", Pattern: "["}}
	_, err := labeler.Label(env, rules)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}
