package classifier_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/classifier"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestClassify_DatabaseKey(t *testing.T) {
	env := makeEnv("DATABASE_URL", "postgres://localhost/db")
	report := classifier.Classify(env)
	if len(report.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(report.Results))
	}
	if report.Results[0].Category != classifier.CategoryDatabase {
		t.Errorf("expected database, got %s", report.Results[0].Category)
	}
}

func TestClassify_AuthKey(t *testing.T) {
	env := makeEnv("JWT_SECRET", "abc123")
	report := classifier.Classify(env)
	if report.Results[0].Category != classifier.CategoryAuth {
		t.Errorf("expected auth, got %s", report.Results[0].Category)
	}
	if report.Results[0].Confidence < 0.9 {
		t.Errorf("expected high confidence, got %.2f", report.Results[0].Confidence)
	}
}

func TestClassify_NetworkKey(t *testing.T) {
	env := makeEnv("API_HOST", "example.com", "API_PORT", "8080")
	report := classifier.Classify(env)
	for _, r := range report.Results {
		if r.Category != classifier.CategoryNetwork {
			t.Errorf("key %s: expected network, got %s", r.Key, r.Category)
		}
	}
}

func TestClassify_ObservabilityKey(t *testing.T) {
	env := makeEnv("OTEL_ENDPOINT", "http://collector:4317")
	report := classifier.Classify(env)
	if report.Results[0].Category != classifier.CategoryObservability {
		t.Errorf("expected observability, got %s", report.Results[0].Category)
	}
}

func TestClassify_StorageKey(t *testing.T) {
	env := makeEnv("S3_BUCKET", "my-bucket")
	report := classifier.Classify(env)
	if report.Results[0].Category != classifier.CategoryStorage {
		t.Errorf("expected storage, got %s", report.Results[0].Category)
	}
}

func TestClassify_GeneralKey(t *testing.T) {
	env := makeEnv("APP_ENV", "production")
	report := classifier.Classify(env)
	if report.Results[0].Category != classifier.CategoryGeneral {
		t.Errorf("expected general, got %s", report.Results[0].Category)
	}
}

func TestClassify_ByCategory_Populated(t *testing.T) {
	env := makeEnv(
		"DB_HOST", "localhost",
		"API_TOKEN", "secret",
		"APP_NAME", "envlens",
	)
	report := classifier.Classify(env)
	if len(report.ByCategory[classifier.CategoryAuth]) == 0 {
		t.Error("expected auth category to be populated")
	}
	if len(report.ByCategory[classifier.CategoryGeneral]) == 0 {
		t.Error("expected general category to be populated")
	}
}

func TestClassify_EmptyEnv(t *testing.T) {
	report := classifier.Classify(map[string]string{})
	if len(report.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(report.Results))
	}
}

func TestClassify_SortedResults(t *testing.T) {
	env := makeEnv("Z_KEY", "z", "A_KEY", "a", "M_KEY", "m")
	report := classifier.Classify(env)
	for i := 1; i < len(report.Results); i++ {
		if report.Results[i-1].Key > report.Results[i].Key {
			t.Errorf("results not sorted: %s > %s", report.Results[i-1].Key, report.Results[i].Key)
		}
	}
}
