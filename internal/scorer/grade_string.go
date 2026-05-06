package scorer

import "fmt"

// String returns a human-readable summary of a Result.
func (r Result) String() string {
	return fmt.Sprintf("Score: %d/100  Grade: %s  Issues: %d", r.Score, r.Grade, len(r.Deductions))
}

// gradeLabel returns a descriptive label for a grade.
func gradeLabel(g Grade) string {
	switch g {
	case GradeA:
		return "Excellent"
	case GradeB:
		return "Good"
	case GradeC:
		return "Fair"
	case GradeD:
		return "Poor"
	default:
		return "Critical"
	}
}

// Summary returns a one-line human-readable summary including the grade label.
func (r Result) Summary() string {
	return fmt.Sprintf("%s — %s (%d/100)", r.Grade, gradeLabel(r.Grade), r.Score)
}
