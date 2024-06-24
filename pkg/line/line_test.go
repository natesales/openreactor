package line

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinearParse(t *testing.T) {
	for _, tc := range []struct {
		in  string
		out Linear
	}{
		{"2x+3", Linear{2, 3}},
		{"-2x+3", Linear{-2, 3}},
		{"2x-3", Linear{2, -3}},
		{"x", Linear{1, 0}},
		{"-x", Linear{-1, 0}},
		{"12x", Linear{12, 0}},
	} {
		t.Run(tc.in, func(t *testing.T) {
			var l Linear
			assert.Nil(t, l.Parse(tc.in))
			assert.Equal(t, tc.out, l)
		})
	}
}

func TestLinearEval(t *testing.T) {
	l := Linear{-2, 3}
	for x, y := range map[float64]float64{
		0: 3,
		1: 1,
		2: -1,
	} {
		assert.Equal(t, y, l.Eval(x))
	}
}
