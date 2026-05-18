package pipeline_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/pipeline"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestPipeline_SingleStage(t *testing.T) {
	p := pipeline.New().Add(pipeline.Stage{
		Name: "uppercase",
		Apply: func(env map[string]string) (map[string]string, error) {
			out := make(map[string]string)
			for k, v := range env {
				out[k] = strings.ToUpper(v)
			}
			return out, nil
		},
	})
	results, err := p.Run(makeEnv("KEY", "hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := pipeline.Final(results)["KEY"]; got != "HELLO" {
		t.Errorf("expected HELLO, got %s", got)
	}
}

func TestPipeline_MultipleStages_Chained(t *testing.T) {
	p := pipeline.New().
		Add(pipeline.Stage{
			Name: "prefix",
			Apply: func(env map[string]string) (map[string]string, error) {
				out := make(map[string]string)
				for k, v := range env {
					out["PRE_"+k] = v
				}
				return out, nil
			},
		}).
		Add(pipeline.Stage{
			Name: "append",
			Apply: func(env map[string]string) (map[string]string, error) {
				out := make(map[string]string)
				for k, v := range env {
					out[k] = v + "_end"
				}
				return out, nil
			},
		})
	results, err := p.Run(makeEnv("FOO", "bar"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	final := pipeline.Final(results)
	if v, ok := final["PRE_FOO"]; !ok || v != "bar_end" {
		t.Errorf("expected PRE_FOO=bar_end, got %v", final)
	}
}

func TestPipeline_StageError_HaltsExecution(t *testing.T) {
	called := false
	p := pipeline.New().
		Add(pipeline.Stage{
			Name: "fail",
			Apply: func(env map[string]string) (map[string]string, error) {
				return nil, errors.New("boom")
			},
		}).
		Add(pipeline.Stage{
			Name: "never",
			Apply: func(env map[string]string) (map[string]string, error) {
				called = true
				return env, nil
			},
		})
	_, err := p.Run(makeEnv("A", "1"))
	if err == nil {
		t.Fatal("expected error")
	}
	if called {
		t.Error("second stage should not have been called")
	}
}

func TestPipeline_ResultsCountMatchesStages(t *testing.T) {
	p := pipeline.New().
		Add(pipeline.Stage{Name: "a", Apply: func(e map[string]string) (map[string]string, error) { return e, nil }}).
		Add(pipeline.Stage{Name: "b", Apply: func(e map[string]string) (map[string]string, error) { return e, nil }})
	results, _ := p.Run(makeEnv())
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestFinal_EmptyResults(t *testing.T) {
	out := pipeline.Final(nil)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
