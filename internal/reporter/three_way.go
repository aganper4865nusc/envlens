package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/envlens/internal/differ"
)

// WriteThreeWay writes a three-way diff report to w in the given format.
func WriteThreeWay(w io.Writer, results []differ.ThreeWayResult, format string) error {
	switch format {
	case "json":
		return writeThreeWayJSON(w, results)
	default:
		return writeThreeWayText(w, results)
	}
}

func writeThreeWayText(w io.Writer, results []differ.ThreeWayResult) error {
	conflicts := 0
	for _, r := range results {
		if r.Conflict {
			conflicts++
		}
	}

	fmt.Fprintf(w, "Three-Way Diff: %d keys, %d conflict(s)\n", len(results), conflicts)
	fmt.Fprintln(w)

	for _, r := range results {
		if r.Conflict {
			fmt.Fprintf(w, "[CONFLICT] %s\n", r.Key)
			fmt.Fprintf(w, "  base:  %s\n", r.Base)
			fmt.Fprintf(w, "  left:  %s\n", r.Left)
			fmt.Fprintf(w, "  right: %s\n", r.Right)
			if r.Resolution != "" {
				fmt.Fprintf(w, "  resolved: %s\n", r.Resolution)
			}
		} else if r.Left != r.Base || r.Right != r.Base {
			fmt.Fprintf(w, "[CHANGED]  %s = %s\n", r.Key, r.Resolution)
		} else {
			fmt.Fprintf(w, "[OK]       %s\n", r.Key)
		}
	}
	return nil
}

type threeWayJSON struct {
	TotalKeys int                   `json:"total_keys"`
	Conflicts int                   `json:"conflicts"`
	Entries   []threeWayEntryJSON   `json:"entries"`
}

type threeWayEntryJSON struct {
	Key        string `json:"key"`
	Base       string `json:"base"`
	Left       string `json:"left"`
	Right      string `json:"right"`
	Conflict   bool   `json:"conflict"`
	Resolution string `json:"resolution,omitempty"`
}

func writeThreeWayJSON(w io.Writer, results []differ.ThreeWayResult) error {
	entries := make([]threeWayEntryJSON, len(results))
	conflicts := 0
	for i, r := range results {
		if r.Conflict {
			conflicts++
		}
		entries[i] = threeWayEntryJSON{
			Key:        r.Key,
			Base:       r.Base,
			Left:       r.Left,
			Right:      r.Right,
			Conflict:   r.Conflict,
			Resolution: r.Resolution,
		}
	}
	payload := threeWayJSON{
		TotalKeys: len(results),
		Conflicts: conflicts,
		Entries:   entries,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
