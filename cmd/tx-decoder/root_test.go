package txdecoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	runCmd := NewRootCommand()
	assert.NotNil(t, runCmd, "run cmd should not be nil")
}
