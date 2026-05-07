package patcher

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadOps reads a patch file and returns a slice of Ops.
//
// Patch file format (one op per line):
//
//	set KEY=VALUE
//	unset KEY
//	rename OLD_KEY NEW_KEY
//	# comment lines and blank lines are ignored
func LoadOps(path string) ([]Op, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("patcher: open %s: %w", path, err)
	}
	defer f.Close()

	var ops []Op
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		op, err := parseLine(line, lineNum)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("patcher: scan %s: %w", path, err)
	}

	return ops, nil
}

// parseLine parses a single non-empty, non-comment line into an Op.
func parseLine(line string, lineNum int) (Op, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return Op{}, fmt.Errorf("patcher: line %d: too few fields: %q", lineNum, line)
	}

	switch strings.ToLower(parts[0]) {
	case "set":
		kv := strings.SplitN(parts[1], "=", 2)
		if len(kv) != 2 {
			return Op{}, fmt.Errorf("patcher: line %d: set requires KEY=VALUE, got %q", lineNum, parts[1])
		}
		return Op{Type: OpSet, Key: kv[0], Value: kv[1]}, nil

	case "unset":
		return Op{Type: OpUnset, Key: parts[1]}, nil

	case "rename":
		if len(parts) < 3 {
			return Op{}, fmt.Errorf("patcher: line %d: rename requires OLD NEW", lineNum)
		}
		return Op{Type: OpRename, Key: parts[1], To: parts[2]}, nil

	default:
		return Op{}, fmt.Errorf("patcher: line %d: unknown directive %q", lineNum, parts[0])
	}
}
