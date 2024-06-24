package line

import (
	"strings"
)

type Line interface {
	// Parse parses a line from a string
	Parse(s string) error

	// Eval evaluates the line at x
	Eval(x float64) float64

	// String returns the string representation of the line
	String() string

	// UnmarshalYAML unmarshalls a line from a YAML node
	UnmarshalYAML(func(any) error) error
}

var (
	_ Line = (*Linear)(nil)
)

func normalize(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", ""))
}

func unmarshal(l Line, u func(any) error) error {
	var s string
	if err := u(&s); err != nil {
		return err
	}
	return l.Parse(s)
}
