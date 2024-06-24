package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSMSequence(t *testing.T) {
	Reset()
	for i := 0; i < len(States); i++ {
		assert.Equal(t, States[i], Get())
		Next()
	}
	Next()
	assert.Equal(t, States[len(States)-1], Get())
}
