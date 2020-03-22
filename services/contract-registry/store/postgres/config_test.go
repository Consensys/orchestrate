package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.NotNil(t, cfg)
}
