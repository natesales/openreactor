package profile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfileParse(t *testing.T) {
	p := `
name: Test Profile
revision: 1
`

	profile, err := Parse([]byte(p))
	assert.Nil(t, err)

	assert.Equal(t, "Test Profile", profile.Name)
	assert.Equal(t, "1", profile.Revision)
	assert.Equal(t, 90000, profile.Vacuum.RotorSpeed)
	assert.Equal(t, "15s", profile.Vacuum.RotorStartupHold.String())
}
