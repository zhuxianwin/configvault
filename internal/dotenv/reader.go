package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Reader reads key-value pairs from a dotenv file.
type Reader struct {
	path string
}

// NewReader creates a new Reader for the given file path.
func NewReader(path string) *Reader {
	return &Reader{path: path}
}

// Read parses the dotenv file and returns a map of key-value pairs.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are also ignored.
func (r *Reader) Read() (map[string]string, error) {
	f, err := os.Open(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("dotenv reader: open file: %w", err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("dotenv reader: invalid format at line %d: %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = stripQuotes(value)

		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("dotenv reader: scan file: %w", err)
	}

	return result, nil
}

// stripQuotes removes surrounding double or single quotes from a value.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
