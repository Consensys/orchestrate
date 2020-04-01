// +build unit

package httpclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client, "Client should not be nil")
}
