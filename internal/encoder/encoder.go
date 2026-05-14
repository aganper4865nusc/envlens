// Package encoder provides utilities for encoding and decoding environment
// variable maps to and from common wire formats (base64, hex, URL-encoded).
package encoder

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
)

// Format represents the encoding format to apply to env values.
type Format string

const (
	FormatBase64 Format = "base64"
	FormatHex    Format = "hex"
	FormatURL    Format = "url"
)

// Result holds the outcome of an encode or decode operation.
type Result struct {
	// Encoded is the transformed map.
	Encoded map[string]string
	// Skipped contains keys whose values could not be processed.
	Skipped []string
}

// Encode encodes all values in env using the given format.
// Keys whose values fail encoding are collected in Result.Skipped.
func Encode(env map[string]string, format Format) (Result, error) {
	result := Result{
		Encoded: make(map[string]string, len(env)),
	}

	keys := sortedKeys(env)
	for _, k := range keys {
		v := env[k]
		encoded, err := encodeValue(v, format)
		if err != nil {
			result.Skipped = append(result.Skipped, k)
			result.Encoded[k] = v
			continue
		}
		result.Encoded[k] = encoded
	}
	return result, nil
}

// Decode decodes all values in env using the given format.
// Keys whose values fail decoding are collected in Result.Skipped.
func Decode(env map[string]string, format Format) (Result, error) {
	result := Result{
		Encoded: make(map[string]string, len(env)),
	}

	keys := sortedKeys(env)
	for _, k := range keys {
		v := env[k]
		decoded, err := decodeValue(v, format)
		if err != nil {
			result.Skipped = append(result.Skipped, k)
			result.Encoded[k] = v
			continue
		}
		result.Encoded[k] = decoded
	}
	return result, nil
}

func encodeValue(v string, format Format) (string, error) {
	switch format {
	case FormatBase64:
		return base64.StdEncoding.EncodeToString([]byte(v)), nil
	case FormatHex:
		return hex.EncodeToString([]byte(v)), nil
	case FormatURL:
		return url.QueryEscape(v), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func decodeValue(v string, format Format) (string, error) {
	switch format {
	case FormatBase64:
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case FormatHex:
		b, err := hex.DecodeString(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case FormatURL:
		decoded, err := url.QueryUnescape(v)
		if err != nil {
			return "", err
		}
		return decoded, nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
