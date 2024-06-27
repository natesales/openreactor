package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMFCCRC(t *testing.T) {
	assert.Equal(t, []byte{0xbe, 0x35}, cksum("foobar"))
	assert.Equal(t, []byte{0x1f, 0xc6}, cksum("test"))
}
