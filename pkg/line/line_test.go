package line

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLinear_Parse(t *testing.T) {
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

func TestLinear_Eval(t *testing.T) {
	l := Linear{-2, 3}
	for x, y := range map[float64]float64{
		0: 3,
		1: 1,
		2: -1,
	} {
		assert.Equal(t, y, l.Eval(x))
	}
}

func TestLinear_UnmarshalYAML(t *testing.T) {
	var l Linear
	assert.Nil(t, yaml.Unmarshal([]byte("2x+3"), &l))
	assert.Equal(t, Linear{2, 3}, l)
}
