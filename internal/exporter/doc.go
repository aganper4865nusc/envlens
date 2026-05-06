// Package exporter provides functionality to serialize parsed environment
// variable maps into multiple output formats suitable for different toolchains.
//
// Supported formats:
//
//   - shell  — POSIX shell export statements (export KEY="value")
//   - docker — Docker-compatible env-file format (KEY=value)
//   - json   — Pretty-printed JSON object
//
// All formats emit keys in lexicographic order for reproducible output.
package exporter
