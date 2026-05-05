package auditor

import (
	"testing"
)

func TestAudit_EmptyValue(t *testing.T) {
	env := map[string]string{"DB_HOST": "", "APP_NAME": "myapp"}
	opts := AuditOptions{FlagEmptyValues: true}
	results := Audit(env, nil, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", results[0].Key)
	}
	if results[0].Severity != "warn" {
		t.Errorf("expected warn severity, got %s", results[0].Severity)
	}
}

func TestAudit_PlainSecret(t *testing.T) {
	env := map[string]string{"API_KEY": "mysupersecret", "APP_NAME": "myapp"}
	opts := AuditOptions{FlagPlainSecrets: true}
	results := Audit(env, nil, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", results[0].Key)
	}
}

func TestAudit_EncodedSecretNotFlagged(t *testing.T) {
	// Value looks encoded (long, mixed case + digits) — should not be flagged.
	env := map[string]string{"API_KEY": "aB3dEfGhIjKlMnOp"}
	opts := AuditOptions{FlagPlainSecrets: true}
	results := Audit(env, nil, opts)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestAudit_DuplicateKeys(t *testing.T) {
	rawLines := []string{
		"DB_HOST=localhost",
		"DB_PORT=5432",
		"DB_HOST=remotehost",
	}
	env := map[string]string{"DB_HOST": "remotehost", "DB_PORT": "5432"}
	opts := AuditOptions{FlagDuplicateKeys: true}
	results := Audit(env, rawLines, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", results[0].Key)
	}
}

func TestAudit_NoDuplicates(t *testing.T) {
	rawLines := []string{"DB_HOST=localhost", "DB_PORT=5432"}
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	opts := AuditOptions{FlagDuplicateKeys: true}
	results := Audit(env, rawLines, opts)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestAudit_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "plainpass",
		"APP_ENV":     "",
	}
	opts := AuditOptions{FlagPlainSecrets: true, FlagEmptyValues: true}
	results := Audit(env, nil, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// Results are sorted by key: APP_ENV < DB_PASSWORD
	if results[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV first, got %s", results[0].Key)
	}
	if results[1].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD second, got %s", results[1].Key)
	}
}
