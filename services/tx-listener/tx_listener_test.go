// +build unit

package txlistener

import (
	"context"
	"testing"
)

type MockSessionManager struct{}

func (m *MockSessionManager) Start(_ context.Context) {}

func TestStart(t *testing.T) {
	txListener := &TxListener{
		manager: &MockSessionManager{},
	}

	txListener.Start(context.Background())
}
