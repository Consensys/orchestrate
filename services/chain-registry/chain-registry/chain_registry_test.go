// +build unit

package chainregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	b := NewBuilder(nil, nil)
	assert.NotNil(t, b, "Builder should not be nil")
}
