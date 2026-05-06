package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envlens/internal/auditor"
	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/linter"
	"github.com/user/envlens/internal/validator"
)

// Report holds all results to be reported.
type Report struct {
	Diffs      []differ.DiffEntry
	Validation []validator.Issue
	Audit      []auditor.Issue
	Lint       []linter.Issue
}

// Write outputs the report to w in the given format ("text" or "json").
func Write(w io.Writer, r Report, format string) error {
	switch format {
	case "json":
		return writeJSON(w, r)
	default:
		return writeText(w, r)
	}
}

func writeText(w io.Writer, r Report) error {
	if len(r.Diffs) == 0 && len(r.Validation) == 0 && len(r.Audit) == 0 && len(r.Lint) == 0 {
		_, err := fmt.Fprintln(w, "No issues found.")
		return err
	}
	for _, d := range r.Diffs {
		_, err := fmt.Fprintf(w, "[diff] %s %s\n", d.Type, d.Key)
		if err != nil {
			return err
		}
	}
	for _, v := range r.Validation {
		_, err := fmt.Fprintf(w, "[validation] %s: %s\n", v.Key, v.Message)
		if err != nil {
			return err
		}
	}
	for _, a := range r.Audit {
		_, err := fmt.Fprintf(w, "[audit] %s: %s\n", a.Key, a.Message)
		if err != nil {
			return err
		}
	}
	for _, l := range r.Lint {
		_, err := fmt.Fprintf(w, "[lint/%s] %s: %s\n", l.Severity, l.Key, l.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, r Report) error {
	type jsonReport struct {
		Diffs      interface{} `json:"diffs"`
		Validation interface{} `json:"validation"`
		Audit      interface{} `json:"audit"`
		Lint       interface{} `json:"lint"`
	}

	diffs := marshalDiff(r.Diffs)
	validation := marshalValidation(r.Validation)
	audit := marshalAudit(r.Audit)
	lint := marshalLint(r.Lint)

	out := jsonReport{Diffs: diffs, Validation: validation, Audit: audit, Lint: lint}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func marshalDiff(entries []differ.DiffEntry) []map[string]string {
	result := make([]map[string]string, 0, len(entries))
	for _, e := range entries {
		result = append(result, map[string]string{"type": string(e.Type), "key": e.Key, "base": e.BaseValue, "target": e.TargetValue})
	}
	return result
}

func marshalValidation(issues []validator.Issue) []map[string]string {
	result := make([]map[string]string, 0, len(issues))
	for _, i := range issues {
		result = append(result, map[string]string{"key": i.Key, "message": i.Message})
	}
	return result
}

func marshalAudit(issues []auditor.Issue) []map[string]string {
	result := make([]map[string]string, 0, len(issues))
	for _, i := range issues {
		result = append(result, map[string]string{"key": i.Key, "message": i.Message})
	}
	return result
}

func marshalLint(issues []linter.Issue) []map[string]string {
	sorted := make([]linter.Issue, len(issues))
	copy(sorted, issues)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })
	result := make([]map[string]string, 0, len(sorted))
	for _, i := range sorted {
		result = append(result, map[string]string{"key": i.Key, "message": i.Message, "severity": i.Severity})
	}
	return result
}
