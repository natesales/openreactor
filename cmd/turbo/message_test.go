package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTurboZeroPad(t *testing.T) {
	assert.Equal(t, "000", zeroPad(0, 3))
	assert.Equal(t, "001", zeroPad(1, 3))
	assert.Equal(t, "001", zeroPad("1", 3))
	assert.Equal(t, "123", zeroPad("123", 3))
}

func TestTurboCksum(t *testing.T) {
	tcs := []struct {
		s string
		e string
	}{
		{"0010000002=?", "095"},
		{"012345", "047"},
	}
	for _, tc := range tcs {
		t.Run(tc.s, func(t *testing.T) {
			assert.Equal(t, tc.e, cksum(tc.s))
		})
	}
}
