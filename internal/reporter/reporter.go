package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/yourorg/envlens/internal/auditor"
	"github.com/yourorg/envlens/internal/differ"
	"github.com/yourorg/envlens/internal/resolver"
	"github.com/yourorg/envlens/internal/validator"
)

// Report bundles all analysis results.
type Report struct {
	Diffs       []differ.Entry
	Validations []validator.Issue
	Audits      []auditor.Finding
	Resolutions []resolver.Resolution
}

// Write serialises the report to w in the requested format ("text" or "json").
func Write(w io.Writer, r Report, format string) error {
	switch format {
	case "json":
		return writeJSON(w, r)
	default:
		return writeText(w, r)
	}
}

func writeText(w io.Writer, r Report) error {
	if len(r.Diffs) == 0 && len(r.Validations) == 0 && len(r.Audits) == 0 {
		_, err := fmt.Fprintln(w, "No issues found.")
		return err
	}
	for _, d := range r.Diffs {
		_, err := fmt.Fprintf(w, "[diff] %s %s\n", d.Key, d.Status)
		if err != nil {
			return err
		}
	}
	for _, v := range r.Validations {
		_, err := fmt.Fprintf(w, "[validation] %s: %s\n", v.Key, v.Message)
		if err != nil {
			return err
		}
	}
	for _, a := range r.Audits {
		_, err := fmt.Fprintf(w, "[audit] %s: %s\n", a.Key, a.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

type jsonReport struct {
	Diffs       []marshalledDiff       `json:"diffs"`
	Validations []marshalledValidation `json:"validations"`
	Audits      []marshalledAudit      `json:"audits"`
	Resolutions []marshalledResolution `json:"resolutions,omitempty"`
}

type marshalledDiff struct {
	Key    string `json:"key"`
	Status string `json:"status"`
}

type marshalledValidation struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

type marshalledAudit struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

type marshalledResolution struct {
	Key      string `json:"key"`
	Source   string `json:"source"`
	Override bool   `json:"override,omitempty"`
}

func writeJSON(w io.Writer, r Report) error {
	out := jsonReport{}

	sort.Slice(r.Diffs, func(i, j int) bool { return r.Diffs[i].Key < r.Diffs[j].Key })
	for _, d := range r.Diffs {
		out.Diffs = append(out.Diffs, marshalledDiff{Key: d.Key, Status: string(d.Status)})
	}
	for _, v := range r.Validations {
		out.Validations = append(out.Validations, marshalledValidation{Key: v.Key, Message: v.Message})
	}
	for _, a := range r.Audits {
		out.Audits = append(out.Audits, marshalledAudit{Key: a.Key, Message: a.Message})
	}
	for _, res := range r.Resolutions {
		out.Resolutions = append(out.Resolutions, marshalledResolution{
			Key:      res.Key,
			Source:   res.Source,
			Override: res.Override,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func marshalDiff(d differ.Entry) marshalledDiff       { return marshalledDiff{Key: d.Key, Status: string(d.Status)} }
func marshalValidation(v validator.Issue) marshalledValidation {
	return marshalledValidation{Key: v.Key, Message: v.Message}
}
func marshalAuditFinding(a auditor.Finding) marshalledAudit {
	return marshalledAudit{Key: a.Key, Message: a.Message}
}
