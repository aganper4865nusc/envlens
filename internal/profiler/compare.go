package profiler

import "fmt"

// Comparison holds a side-by-side diff of two profiles.
type Comparison struct {
	SourceA       string
	SourceB       string
	KeyDelta      int // B.TotalKeys - A.TotalKeys
	EmptyDelta    int
	NewSensitive  []string // sensitive keys in B not in A
	LostSensitive []string // sensitive keys in A not in B
	Notes         []string
}

// Compare produces a Comparison between two Profiles.
func Compare(a, b Profile) Comparison {
	c := Comparison{
		SourceA:    a.Source,
		SourceB:    b.Source,
		KeyDelta:   b.TotalKeys - a.TotalKeys,
		EmptyDelta: b.EmptyValues - a.EmptyValues,
	}

	aSet := toSet(a.SensitiveKeys)
	bSet := toSet(b.SensitiveKeys)

	for k := range bSet {
		if _, ok := aSet[k]; !ok {
			c.NewSensitive = append(c.NewSensitive, k)
		}
	}
	for k := range aSet {
		if _, ok := bSet[k]; !ok {
			c.LostSensitive = append(c.LostSensitive, k)
		}
	}

	if c.KeyDelta > 0 {
		c.Notes = append(c.Notes, fmt.Sprintf("%s has %d more keys than %s", b.Source, c.KeyDelta, a.Source))
	} else if c.KeyDelta < 0 {
		c.Notes = append(c.Notes, fmt.Sprintf("%s has %d fewer keys than %s", b.Source, -c.KeyDelta, a.Source))
	}
	if c.EmptyDelta > 0 {
		c.Notes = append(c.Notes, fmt.Sprintf("%s has more empty values than %s", b.Source, a.Source))
	}

	return c
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
