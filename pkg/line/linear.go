package line

import (
	"fmt"
	"strings"
)

// Linear represents a linear equation in slope intercept form
type Linear struct {
	m, b float64
}

func (l *Linear) Parse(s string) error {
	s = normalize(s)

	var err error
	if s == "x" {
		return l.Parse("1x")
	} else if s == "-x" {
		return l.Parse("-1x")
	} else if strings.Contains(s, "x-") {
		_, err = fmt.Sscanf(s, "%fx-%f", &l.m, &l.b)
		l.b = -l.b
	} else if strings.Contains(s, "x+") {
		_, err = fmt.Sscanf(s, "%fx+%f", &l.m, &l.b)
	} else {
		l.b = 0
		_, err = fmt.Sscanf(s, "%fx", &l.m)
	}

	if err != nil {
		return err
	}
	return nil
}

func (l *Linear) Eval(x float64) float64 {
	return l.m*x + l.b
}

func (l *Linear) String() string {
	return fmt.Sprintf("%fx+%f", l.m, l.b)
}

func (l *Linear) UnmarshalYAML(u func(any) error) error {
	return unmarshal(l, u)
}
