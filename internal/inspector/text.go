package inspector

import (
	"fmt"
	"io"
	"strings"
)

// WriteText writes a human-readable inspection report to w.
func WriteText(w io.Writer, r Report) {
	fmt.Fprintf(w, "=== Inspection: %s ===\n", r.Source)
	fmt.Fprintf(w, "Keys: %d  Empty: %d  Sensitive: %d\n",
		r.Profile.TotalKeys, r.Profile.EmptyValues, r.Profile.SensitiveKeys)
	fmt.Fprintf(w, "Score: %d (%s)\n", r.Score.Value, r.Score.Grade)

	if len(r.Lint) > 0 {
		fmt.Fprintln(w, "\nLint Issues:")
		for _, issue := range r.Lint {
			fmt.Fprintf(w, "  [%s] %s: %s\n",
				strings.ToUpper(string(issue.Severity)), issue.Key, issue.Message)
		}
	}

	if len(r.Audit) > 0 {
		fmt.Fprintln(w, "\nAudit Issues:")
		for _, issue := range r.Audit {
			fmt.Fprintf(w, "  [%s] %s: %s\n",
				strings.ToUpper(issue.Severity), issue.Key, issue.Message)
		}
	}

	if !HasIssues(r) {
		fmt.Fprintln(w, "No issues found.")
	}
}
