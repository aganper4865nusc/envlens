// Package pipeline provides a composable, sequential env transformation pipeline
// that chains multiple named stages, collecting results and errors per stage.
package pipeline

import "fmt"

// Stage represents a single named transformation step in a pipeline.
type Stage struct {
	Name    string
	Apply   func(env map[string]string) (map[string]string, error)
}

// Result holds the output of a single stage execution.
type Result struct {
	Stage  string
	Output map[string]string
	Err    error
}

// Pipeline is an ordered sequence of stages.
type Pipeline struct {
	stages []Stage
}

// New creates an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Add appends a stage to the pipeline.
func (p *Pipeline) Add(s Stage) *Pipeline {
	p.stages = append(p.stages, s)
	return p
}

// Run executes all stages in order, passing each output as the next input.
// Execution halts on the first error. All per-stage results are returned.
func (p *Pipeline) Run(initial map[string]string) ([]Result, error) {
	results := make([]Result, 0, len(p.stages))
	current := copyMap(initial)

	for _, stage := range p.stages {
		out, err := stage.Apply(current)
		results = append(results, Result{Stage: stage.Name, Output: out, Err: err})
		if err != nil {
			return results, fmt.Errorf("pipeline stage %q failed: %w", stage.Name, err)
		}
		current = out
	}
	return results, nil
}

// Final returns the output of the last successful stage, or the initial map
// if no stages have run.
func Final(results []Result) map[string]string {
	for i := len(results) - 1; i >= 0; i-- {
		if results[i].Err == nil && results[i].Output != nil {
			return results[i].Output
		}
	}
	return map[string]string{}
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
