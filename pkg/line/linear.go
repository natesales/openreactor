package line

import (
	"fmt"
	"strconv"
	"strings"
)

// Linear represents a linear equation in slope intercept form
type Linear struct {
	m, b float64
}

// FromSlopeIntercept creates a new linear equation from a slope and y-intercept
func FromSlopeIntercept(m, b float64) *Linear {
	return &Linear{m: m, b: b}
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
	if l.m == 0 && l.b == 0 {
		return "0"
	}

	return fmt.Sprintf(
		"%sx+%s",
		strconv.FormatFloat(l.m, 'f', -1, 64),
		strconv.FormatFloat(l.b, 'f', -1, 64),
	)
}

func (l *Linear) UnmarshalYAML(u func(any) error) error {
	return unmarshal(l, u)
}

func (l *Linear) MarshalYAML() (interface{}, error) {
	return marshal(l)
}
