package pipeline

import (
	"strings"

	"github.com/yourorg/envlens/internal/transformer"
)

// UppercaseKeysStage returns a Stage that uppercases all env keys.
func UppercaseKeysStage() Stage {
	return Stage{
		Name: "uppercase-keys",
		Apply: func(env map[string]string) (map[string]string, error) {
			return transformer.Transform(env, transformer.Options{UppercaseKeys: true}), nil
		},
	}
}

// AddPrefixStage returns a Stage that adds a prefix to every key.
func AddPrefixStage(prefix string) Stage {
	return Stage{
		Name: "add-prefix:" + prefix,
		Apply: func(env map[string]string) (map[string]string, error) {
			return transformer.Transform(env, transformer.Options{AddPrefix: prefix}), nil
		},
	}
}

// StripPrefixStage returns a Stage that removes a prefix from every key.
func StripPrefixStage(prefix string) Stage {
	return Stage{
		Name: "strip-prefix:" + prefix,
		Apply: func(env map[string]string) (map[string]string, error) {
			return transformer.Transform(env, transformer.Options{StripPrefix: prefix}), nil
		},
	}
}

// DropEmptyStage returns a Stage that removes keys with empty values.
func DropEmptyStage() Stage {
	return Stage{
		Name: "drop-empty",
		Apply: func(env map[string]string) (map[string]string, error) {
			out := make(map[string]string)
			for k, v := range env {
				if strings.TrimSpace(v) != "" {
					out[k] = v
				}
			}
			return out, nil
		},
	}
}
