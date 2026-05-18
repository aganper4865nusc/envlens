// Package pipeline provides a composable sequential execution model for
// environment variable transformations.
//
// A Pipeline is built by adding named Stage values, each of which accepts a
// map[string]string and returns a transformed copy (or an error). Stages are
// executed in insertion order; the output of each stage becomes the input of
// the next. Execution halts on the first error.
//
// Pre-built stages for common operations (uppercase keys, add/strip prefix,
// drop empty values) are available in builder.go and integrate with the
// internal/transformer package.
//
// Example:
//
//	p := pipeline.New().
//		Add(pipeline.DropEmptyStage()).
//		Add(pipeline.UppercaseKeysStage()).
//		Add(pipeline.AddPrefixStage("APP_"))
//
//	results, err := p.Run(env)
//	final := pipeline.Final(results)
package pipeline
