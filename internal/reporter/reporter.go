package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/auditor"
	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/validator"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Report aggregates results from diff, validation, and audit operations.
type Report struct {
	DiffResults      []differ.DiffEntry
	ValidateResults  []validator.ValidationIssue
	AuditResults     []auditor.AuditIssue
}

// Write renders the report to the given writer in the specified format.
func Write(w io.Writer, r Report, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, r)
	default:
		return writeText(w, r)
	}
}

func writeText(w io.Writer, r Report) error {
	if len(r.DiffResults) > 0 {
		fmt.Fprintln(w, "=== Diff Results ===")
		for _, entry := range r.DiffResults {
			switch entry.Status {
			case differ.StatusAdded:
				fmt.Fprintf(w, "  + %s=%s\n", entry.Key, entry.NewValue)
			case differ.StatusRemoved:
				fmt.Fprintf(w, "  - %s=%s\n", entry.Key, entry.OldValue)
			case differ.StatusChanged:
				fmt.Fprintf(w, "  ~ %s: %q -> %q\n", entry.Key, entry.OldValue, entry.NewValue)
			case differ.StatusUnchanged:
				fmt.Fprintf(w, "    %s=%s\n", entry.Key, entry.NewValue)
			}
		}
	}

	if len(r.ValidateResults) > 0 {
		fmt.Fprintln(w, "=== Validation Issues ===")
		for _, issue := range r.ValidateResults {
			fmt.Fprintf(w, "  [%s] %s: %s\n", strings.ToUpper(issue.Severity), issue.Key, issue.Message)
		}
	}

	if len(r.AuditResults) > 0 {
		fmt.Fprintln(w, "=== Audit Issues ===")
		for _, issue := range r.AuditResults {
			fmt.Fprintf(w, "  [%s] %s: %s\n", strings.ToUpper(issue.Severity), issue.Key, issue.Message)
		}
	}

	if len(r.DiffResults) == 0 && len(r.ValidateResults) == 0 && len(r.AuditResults) == 0 {
		fmt.Fprintln(w, "No issues found.")
	}

	return nil
}

func writeJSON(w io.Writer, r Report) error {
	fmt.Fprintf(w, "{\n")
	fmt.Fprintf(w, "  \"diff\": %s,\n", marshalDiff(r.DiffResults))
	fmt.Fprintf(w, "  \"validation\": %s,\n", marshalValidation(r.ValidateResults))
	fmt.Fprintf(w, "  \"audit\": %s\n", marshalAudit(r.AuditResults))
	fmt.Fprintf(w, "}\n")
	return nil
}

func marshalDiff(entries []differ.DiffEntry) string {
	if len(entries) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, e := range entries {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, `{"key":%q,"status":%q,"old":%q,"new":%q}`, e.Key, e.Status, e.OldValue, e.NewValue)
	}
	sb.WriteString("]")
	return sb.String()
}

func marshalValidation(issues []validator.ValidationIssue) string {
	if len(issues) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range issues {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, `{"key":%q,"severity":%q,"message":%q}`, v.Key, v.Severity, v.Message)
	}
	sb.WriteString("]")
	return sb.String()
}

func marshalAudit(issues []auditor.AuditIssue) string {
	if len(issues) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, a := range issues {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(&sb, `{"key":%q,"severity":%q,"message":%q}`, a.Key, a.Severity, a.Message)
	}
	sb.WriteString("]")
	return sb.String()
}
