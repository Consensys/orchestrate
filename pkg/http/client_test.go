// +build unit

package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	client := NewClient(NewDefaultConfig())
	assert.NotNil(t, client, "Client should not be nil")
}
