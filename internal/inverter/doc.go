// Package inverter flips an environment variable map so that values become
// keys and keys become values.
//
// This is useful for reverse-lookup scenarios, such as finding which variable
// holds a given endpoint URL or token value. When multiple keys share the same
// value a configurable collision policy (first, last, or skip) determines
// which entry survives in the output.
//
// Basic usage:
//
//	res := inverter.Invert(env, inverter.DefaultOptions())
//	fmt.Println(res.Inverted["https://api.example.com"]) // → API_URL
package inverter
