package envchain_test

import (
	"errors"
	"testing"

	"github.com/user/envlens/internal/envchain"
)

func baseEnv() map[string]string {
	return map[string]string{"A": "1", "B": "2", "C": "3"}
}

func TestChain_SingleStep(t *testing.T) {
	c := envchain.New().Add("uppercase-values", func(env map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(env))
		for k, v := range env {
			out[k] = v + "!"
		}
		return out, nil
	})

	res, err := c.Run(baseEnv())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["A"] != "1!" {
		t.Errorf("expected '1!', got %q", res.Env["A"])
	}
	if len(res.Log) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(res.Log))
	}
}

func TestChain_MultipleSteps_OrderPreserved(t *testing.T) {
	var order []string
	c := envchain.New().
		Add("first", func(env map[string]string) (map[string]string, error) {
			order = append(order, "first")
			return env, nil
		}).
		Add("second", func(env map[string]string) (map[string]string, error) {
			order = append(order, "second")
			return env, nil
		})

	_, err := c.Run(baseEnv())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order[0] != "first" || order[1] != "second" {
		t.Errorf("wrong order: %v", order)
	}
}

func TestChain_StepError_HaltsExecution(t *testing.T) {
	ran := false
	c := envchain.New().
		Add("fail", func(env map[string]string) (map[string]string, error) {
			return nil, errors.New("boom")
		}).
		Add("should-not-run", func(env map[string]string) (map[string]string, error) {
			ran = true
			return env, nil
		})

	_, err := c.Run(baseEnv())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if ran {
		t.Error("second step should not have run after failure")
	}
}

func TestChain_LogRecordsBeforeAfter(t *testing.T) {
	c := envchain.New().Add("drop-one", func(env map[string]string) (map[string]string, error) {
		out := make(map[string]string)
		for k, v := range env {
			if k != "C" {
				out[k] = v
			}
		}
		return out, nil
	})

	res, err := c.Run(baseEnv())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Log[0].Before != 3 {
		t.Errorf("expected before=3, got %d", res.Log[0].Before)
	}
	if res.Log[0].After != 2 {
		t.Errorf("expected after=2, got %d", res.Log[0].After)
	}
}

func TestChain_OriginalInputUnmodified(t *testing.T) {
	input := baseEnv()
	c := envchain.New().Add("mutate", func(env map[string]string) (map[string]string, error) {
		env["NEW"] = "injected"
		return env, nil
	})

	_, _ = c.Run(input)
	if _, ok := input["NEW"]; ok {
		t.Error("original input map was modified")
	}
}
