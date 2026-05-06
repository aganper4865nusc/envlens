// Package scorer evaluates the overall quality of an environment variable file
// by aggregating findings from the auditor, linter, and validator packages into
// a single numeric score (0–100) and a letter grade (A–F).
//
// A perfect file with no issues receives a score of 100 (grade A). Each
// detected problem deducts a weighted number of points depending on its
// category and severity.
package scorer
