package differ

import "sort"

// Annotation holds a diff entry enriched with contextual metadata.
type Annotation struct {
	Key      string
	Status   string // "added", "removed", "changed", "unchanged"
	OldValue string
	NewValue string
	Tags     []string
}

// AnnotateOptions controls how annotations are applied.
type AnnotateOptions struct {
	// TagAdded is the tag applied to keys that are new in the target.
	TagAdded string
	// TagRemoved is the tag applied to keys missing from the target.
	TagRemoved string
	// TagChanged is the tag applied to keys whose value changed.
	TagChanged string
	// TagUnchanged is the tag applied to keys that are identical.
	TagUnchanged string
	// ExtraTags is a static set of tags added to every annotation.
	ExtraTags []string
}

// DefaultAnnotateOptions returns sensible defaults.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{
		TagAdded:     "added",
		TagRemoved:   "removed",
		TagChanged:   "changed",
		TagUnchanged: "unchanged",
	}
}

// Annotate enriches a slice of DiffEntry values with tags and metadata.
func Annotate(entries []DiffEntry, opts AnnotateOptions) []Annotation {
	out := make([]Annotation, 0, len(entries))

	for _, e := range entries {
		a := Annotation{
			Key:      e.Key,
			Status:   e.Status,
			OldValue: e.OldValue,
			NewValue: e.NewValue,
			Tags:     make([]string, 0),
		}

		switch e.Status {
		case "added":
			a.Tags = append(a.Tags, opts.TagAdded)
		case "removed":
			a.Tags = append(a.Tags, opts.TagRemoved)
		case "changed":
			a.Tags = append(a.Tags, opts.TagChanged)
		default:
			a.Tags = append(a.Tags, opts.TagUnchanged)
		}

		a.Tags = append(a.Tags, opts.ExtraTags...)
		sort.Strings(a.Tags)
		out = append(out, a)
	}

	return out
}

// FilterAnnotations returns only annotations whose status matches one of the given statuses.
func FilterAnnotations(annotations []Annotation, statuses ...string) []Annotation {
	set := make(map[string]struct{}, len(statuses))
	for _, s := range statuses {
		set[s] = struct{}{}
	}

	out := make([]Annotation, 0)
	for _, a := range annotations {
		if _, ok := set[a.Status]; ok {
			out = append(out, a)
		}
	}
	return out
}
