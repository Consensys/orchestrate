// +build unit

package loadbalancer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	b := NewBuilder()
	require.NotNil(t, b, "Builder should have been created")
}
