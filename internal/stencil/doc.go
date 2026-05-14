// Package stencil generates template .env files from real environment maps.
//
// It replaces values with typed placeholders, making it easy to produce
// safe, shareable onboarding templates or documentation stubs without
// leaking sensitive credentials.
//
// Sensitive keys (those matching patterns like SECRET, TOKEN, PASSWORD, KEY)
// are always replaced with a descriptive placeholder regardless of options.
// Non-sensitive keys can optionally preserve their original values when
// PreserveDefaults is enabled.
//
// Example usage:
//
//	env := map[string]string{"DB_PASSWORD": "s3cr3t", "APP_PORT": "8080"}
//	opts := stencil.DefaultOptions()
//	entries := stencil.Generate(env, opts)
//	fmt.Print(stencil.Render(entries, opts.CommentPrefix))
package stencil
