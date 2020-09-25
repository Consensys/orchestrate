// +build unit

package txlistener

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockSessionManager struct{}

func (m *MockSessionManager) Run(_ context.Context) error { return nil }

func TestRun(t *testing.T) {
	txListener := &TxListener{
		manager: &MockSessionManager{},
	}

	err := txListener.Run(context.Background())
	assert.NoError(t, err)
}
