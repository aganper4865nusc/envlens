// Package envchain provides a pipeline-style chaining mechanism for applying
// sequential transformations to an environment variable map.
package envchain

import "fmt"

// StepFunc is a function that transforms an env map and returns a new map or an error.
type StepFunc func(env map[string]string) (map[string]string, error)

// Step represents a named transformation step in the chain.
type Step struct {
	Name string
	Fn   StepFunc
}

// Chain holds an ordered list of transformation steps.
type Chain struct {
	steps []Step
}

// Result holds the final env map and a log of each step's outcome.
type Result struct {
	Env  map[string]string
	Log  []StepLog
}

// StepLog records what happened during a single step.
type StepLog struct {
	Step    string
	Before  int
	After   int
	Err     error
}

// New creates an empty Chain.
func New() *Chain {
	return &Chain{}
}

// Add appends a named step to the chain.
func (c *Chain) Add(name string, fn StepFunc) *Chain {
	c.steps = append(c.steps, Step{Name: name, Fn: fn})
	return c
}

// Run executes all steps in order, passing the output of each step as input to the next.
// Execution halts on the first error.
func (c *Chain) Run(initial map[string]string) (*Result, error) {
	current := copyMap(initial)
	logs := make([]StepLog, 0, len(c.steps))

	for _, step := range c.steps {
		before := len(current)
		next, err := step.Fn(current)
		entry := StepLog{
			Step:   step.Name,
			Before: before,
			Err:    err,
		}
		if err != nil {
			entry.After = before
			logs = append(logs, entry)
			return &Result{Env: current, Log: logs}, fmt.Errorf("step %q failed: %w", step.Name, err)
		}
		entry.After = len(next)
		logs = append(logs, entry)
		current = next
	}

	return &Result{Env: current, Log: logs}, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
