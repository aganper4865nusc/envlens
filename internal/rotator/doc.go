// Package rotator applies key rotation operations to an environment map.
//
// A rotation renames a key from one name to another, optionally transforming
// its value in the process. Rotations are useful when migrating configuration
// schemas across deployment stages without losing values.
//
// Usage:
//
//	ops := []rotator.Op{
//		{FromKey: "DB_PASS", ToKey: "DATABASE_PASSWORD"},
//		{FromKey: "API_KEY", ToKey: "SERVICE_API_KEY", Transform: strings.ToUpper},
//	}
//	res, err := rotator.Rotate(env, ops)
package rotator
