package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
)

// WriteAnnotated writes annotated diff results to w in the given format.
func WriteAnnotated(w io.Writer, annotations []differ.Annotation, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return writeAnnotatedJSON(w, annotations)
	default:
		return writeAnnotatedText(w, annotations)
	}
}

func writeAnnotatedText(w io.Writer, annotations []differ.Annotation) error {
	if len(annotations) == 0 {
		_, err := fmt.Fprintln(w, "No diff annotations.")
		return err
	}

	for _, a := range annotations {
		tags := strings.Join(a.Tags, ", ")
		switch a.Status {
		case "added":
			fmt.Fprintf(w, "[+] %s = %q  [%s]\n", a.Key, a.NewValue, tags)
		case "removed":
			fmt.Fprintf(w, "[-] %s = %q  [%s]\n", a.Key, a.OldValue, tags)
		case "changed":
			fmt.Fprintf(w, "[~] %s: %q -> %q  [%s]\n", a.Key, a.OldValue, a.NewValue, tags)
		default:
			fmt.Fprintf(w, "[=] %s = %q  [%s]\n", a.Key, a.NewValue, tags)
		}
	}

	fmt.Fprintf(w, "\nTotal: %d annotation(s)\n", len(annotations))
	return nil
}

type annotatedJSON struct {
	Key      string   `json:"key"`
	Status   string   `json:"status"`
	OldValue string   `json:"old_value,omitempty"`
	NewValue string   `json:"new_value,omitempty"`
	Tags     []string `json:"tags"`
}

func writeAnnotatedJSON(w io.Writer, annotations []differ.Annotation) error {
	records := make([]annotatedJSON, 0, len(annotations))
	for _, a := range annotations {
		records = append(records, annotatedJSON{
			Key:      a.Key,
			Status:   a.Status,
			OldValue: a.OldValue,
			NewValue: a.NewValue,
			Tags:     a.Tags,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}
