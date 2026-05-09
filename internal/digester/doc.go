// Package digester provides deterministic fingerprinting of environment variable maps.
//
// A digest is a SHA-256 hash computed over the sorted key=value pairs of an env map,
// making it suitable for change detection, caching, and deployment auditing.
//
// Example usage:
//
//	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
//	result := digester.Digest(env, digester.Options{})
//	fmt.Println(result.Fingerprint) // deterministic hex string
package digester
