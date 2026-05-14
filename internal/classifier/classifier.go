// Package classifier categorizes environment variable keys into semantic groups
// such as database, network, auth, observability, and general configuration.
package classifier

import (
	"regexp"
	"sort"
	"strings"
)

// Category represents a semantic classification for an env key.
type Category string

const (
	CategoryDatabase      Category = "database"
	CategoryNetwork       Category = "network"
	CategoryAuth          Category = "auth"
	CategoryObservability Category = "observability"
	CategoryStorage       Category = "storage"
	CategoryGeneral       Category = "general"
)

// Result holds the classification result for a single key.
type Result struct {
	Key      string
	Category Category
	Confidence float64 // 0.0–1.0
}

// Report is the full classification output.
type Report struct {
	Results    []Result
	ByCategory map[Category][]string
}

type rule struct {
	pattern    *regexp.Regexp
	category   Category
	confidence float64
}

var defaultRules = []rule{
	{regexp.MustCompile(`(?i)(DB|DATABASE|POSTGRES|MYSQL|MONGO|REDIS|SQLITE|DSN|JDBC)`), CategoryDatabase, 0.9},
	{regexp.MustCompile(`(?i)(HOST|PORT|URL|URI|ADDR|ENDPOINT|PROXY|DOMAIN|DNS)`), CategoryNetwork, 0.85},
	{regexp.MustCompile(`(?i)(SECRET|TOKEN|PASSWORD|PASSWD|API_KEY|APIKEY|AUTH|CREDENTIAL|CERT|PRIVATE)`), CategoryAuth, 0.95},
	{regexp.MustCompile(`(?i)(LOG|TRACE|METRIC|OTEL|SENTRY|DATADOG|NEWRELIC|MONITOR|DEBUG)`), CategoryObservability, 0.85},
	{regexp.MustCompile(`(?i)(S3|BUCKET|BLOB|STORAGE|GCS|MINIO|CDN|UPLOAD)`), CategoryStorage, 0.9},
}

// Classify categorizes each key in the provided env map.
func Classify(env map[string]string) Report {
	report := Report{
		Results:    make([]Result, 0, len(env)),
		ByCategory: make(map[Category][]string),
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		res := classifyKey(key)
		report.Results = append(report.Results, res)
		report.ByCategory[res.Category] = append(report.ByCategory[res.Category], key)
	}

	return report
}

func classifyKey(key string) Result {
	upper := strings.ToUpper(key)
	best := Result{Key: key, Category: CategoryGeneral, Confidence: 0.5}

	for _, r := range defaultRules {
		if r.pattern.MatchString(upper) && r.confidence > best.Confidence {
			best.Category = r.category
			best.Confidence = r.confidence
		}
	}

	return best
}
