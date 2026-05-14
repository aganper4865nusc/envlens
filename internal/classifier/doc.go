// Package classifier provides semantic classification of environment variable keys.
//
// It inspects key names and assigns them to categories such as:
//
//   - database  — keys related to DB connections (DB_URL, POSTGRES_DSN, etc.)
//   - network   — keys related to hosts, ports, and endpoints
//   - auth      — keys related to secrets, tokens, and credentials
//   - observability — keys related to logging, tracing, and metrics
//   - storage   — keys related to object storage (S3, GCS, etc.)
//   - general   — keys that do not match any specific category
//
// Each result includes a confidence score between 0.0 and 1.0 indicating
// how strongly the key matched a known pattern.
//
// Example:
//
//	report := classifier.Classify(env)
//	for _, r := range report.Results {
//		fmt.Printf("%s => %s (%.0f%%)\n", r.Key, r.Category, r.Confidence*100)
//	}
package classifier
