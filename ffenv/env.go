package ffenv

import (
	"bufio"
	"io"
	"strings"
)

// Parser is a parser for config files in an extremely simple format. Each
// line is tokenized as a single key/value pair. The first =-delimited
// token in the line is interpreted as the flag name, and all remaining tokens
// are interpreted as the value. Any leading hyphens on the flag name are
// ignored.
func Parser(r io.Reader, set func(name, value string) error) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue // skip empties
		}

		if line[0] == '#' {
			continue // skip comments
		}

		var (
			name  string
			value string
			index = strings.IndexRune(line, '=')
		)
		if index < 0 {
			name, value = line, "true" // boolean option
		} else {
			name, value = strings.ToLower(line[:index]), line[index+1:]
			name = strings.ReplaceAll(name, "_", "-")
		}

		if i := strings.Index(value, " #"); i >= 0 {
			value = strings.TrimSpace(value[:i])
		}

		if err := set(name, value); err != nil {
			return err
		}
	}
	return nil
}
