// +build unit
// +build !race

package session

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/dynamic"
)

const (
	keyError = "error"
)

type MockSession struct {
	hasRan chan struct{}
}

func NewMockSession() *MockSession {
	return &MockSession{
		hasRan: make(chan struct{}),
	}
}

func (s *MockSession) Run(ctx context.Context) error {
	if ctx.Value(keyError) != nil {
		return fmt.Errorf("test")
	}

	<-ctx.Done()
	// Simulate some latency before finishing
	time.Sleep(50 * time.Millisecond)
	close(s.hasRan)
	return nil
}

type MockBuilder struct {
	mux      *sync.RWMutex
	sessions map[string]*MockSession
}

func (b *MockBuilder) NewSession(chain *dynamic.Chain) (Session, error) {
	if chain.UUID == keyError {
		return nil, fmt.Errorf("error")
	}

	return b.getSession(chain.UUID), nil
}

func (b *MockBuilder) addSession(key string, sess *MockSession) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.sessions[key] = sess
}

func (b *MockBuilder) getSession(key string) *MockSession {
	b.mux.RLock()
	defer b.mux.RUnlock()
	sess, ok := b.sessions[key]
	if !ok {
		panic("no session")
	}
	return sess
}

type MockProvider struct{}

func (p *MockProvider) Run(ctx context.Context, _ chan<- *dynamic.Message) error {
	if ctx.Value(keyError) != nil {
		return fmt.Errorf("test")
	}
	<-ctx.Done()
	return nil
}

func TestManager(t *testing.T) {
	prvdr := &MockProvider{}
	builder := &MockBuilder{
		mux:      &sync.RWMutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, prvdr)

	// Prepare 2 sessions
	testChain1, testChain2 := "test-chain-1", "test-chain-2"
	session1, session2 := NewMockSession(), NewMockSession()
	builder.addSession(testChain1, session1)
	builder.addSession(testChain2, session2)

	// Start manager
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		_ = manager.Run(ctx)
		close(done)
	}()

	// Send command to start both session
	manager.commands <- &Command{
		Type: START,
		Chain: &dynamic.Chain{
			UUID: testChain1,
		},
	}

	manager.commands <- &Command{
		Type: START,
		Chain: &dynamic.Chain{
			UUID: testChain2,
		},
	}

	// Stop first session
	manager.commands <- &Command{
		Type: STOP,
		Chain: &dynamic.Chain{
			UUID: testChain1,
		},
	}

	// Session 1 should have completed
	select {
	case <-time.After(time.Second):
		t.Errorf("Session 1 did not complete")
	case <-session1.hasRan:
	}

	// Cancel to stop Manager
	cancel()

	// Manager should have completed
	select {
	case <-time.After(time.Second):
		t.Errorf("Manager did not complete")
	case <-done:
	}

	// Session 2 should have completed
	select {
	case <-time.After(time.Second):
		t.Errorf("Session 2 did not complete")
	case <-session2.hasRan:
	}
}

func TestErrors(t *testing.T) {
	provider := &MockProvider{}
	builder := &MockBuilder{
		mux:      &sync.RWMutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, provider)
	err := fmt.Errorf("test")
	go func() {
		manager.errors <- err
	}()
}

func TestListenProvider(t *testing.T) {
	provider := &MockProvider{}
	builder := &MockBuilder{
		mux:      &sync.RWMutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, provider)

	ctx := context.WithValue(context.Background(), keyError, true) // nolint
	go func() {
		manager.listenProvider(ctx)
	}()
}

func TestListenConfiguration(t *testing.T) {
	provider := &MockProvider{}
	builder := &MockBuilder{
		mux:      &sync.RWMutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, provider)
	chain1 := &dynamic.Chain{UUID: "test"}

	go func() {
		manager.msgInput <- &dynamic.Message{
			Provider: "test",
			Configuration: &dynamic.Configuration{Chains: map[string]*dynamic.Chain{
				"test": chain1,
			}},
		}
	}()
	go func() { manager.listenConfiguration() }()
	cmd := <-manager.commands
	assert.Equal(t, cmd.Type, START, "should get start command")
	assert.Equal(t, cmd.Chain, chain1, "should get start command")
}

func TestExecuteCommand(t *testing.T) {
	provider := &MockProvider{}
	builder := &MockBuilder{
		mux:      &sync.RWMutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, provider)

	chain := &dynamic.Chain{UUID: "test", TenantID: "test", Name: "test"}
	session1 := NewMockSession()
	builder.addSession(chain.UUID, session1)

	cmd := &Command{
		Type:  UPDATE,
		Chain: chain,
	}
	manager.executeCommand(context.Background(), cmd)
}

func TestRunSession(t *testing.T) {
	provider := &MockProvider{}
	builder := &MockBuilder{
		mux:      &sync.RWMutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, provider)

	chain := &dynamic.Chain{UUID: keyError}
	ctx := context.WithValue(context.Background(), keyError, true) // nolint

	go func() {
		manager.runSession(ctx, chain)
	}()
}
