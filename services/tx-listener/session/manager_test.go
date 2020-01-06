package session

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
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
	mux      *sync.Mutex
	sessions map[string]*MockSession
}

func (b *MockBuilder) NewSession(node *dynamic.Node) (Session, error) {
	if node.ID == keyError {
		return nil, fmt.Errorf("test")
	}

	return b.getSession(node.ID), nil
}

func (b *MockBuilder) addSession(key string, sess *MockSession) {
	b.mux.Lock()
	b.sessions[key] = sess
	b.mux.Unlock()
}

func (b *MockBuilder) getSession(key string) *MockSession {
	b.mux.Lock()
	defer b.mux.Unlock()
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
		mux:      &sync.Mutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, prvdr)

	// Prepare 2 sessions
	testNode1, testNode2 := "test-node-1", "test-node-2"
	session1, session2 := NewMockSession(), NewMockSession()
	builder.addSession(testNode1, session1)
	builder.addSession(testNode2, session2)

	// Start manager
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		manager.Start(ctx)
		close(done)
	}()

	// Send command to start both session
	manager.commands <- &Command{
		Type: START,
		Node: &dynamic.Node{
			ID: testNode1,
		},
	}

	manager.commands <- &Command{
		Type: START,
		Node: &dynamic.Node{
			ID: testNode2,
		},
	}

	// Stop first session
	manager.commands <- &Command{
		Type: STOP,
		Node: &dynamic.Node{
			ID: testNode1,
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
		mux:      &sync.Mutex{},
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
		mux:      &sync.Mutex{},
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
		mux:      &sync.Mutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, provider)
	node1 := &dynamic.Node{ID: "test"}

	go func() {
		manager.msgInput <- &dynamic.Message{
			Provider: "test",
			Configuration: &dynamic.Configuration{Nodes: map[string]*dynamic.Node{
				"test": node1,
			}},
		}
	}()
	go func() { manager.listenConfiguration() }()
	cmd := <-manager.commands
	assert.Equal(t, cmd.Type, START, "should get start command")
	assert.Equal(t, cmd.Node, node1, "should get start command")
}

func TestExecuteCommand(t *testing.T) {
	provider := &MockProvider{}
	builder := &MockBuilder{
		mux:      &sync.Mutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, provider)

	node := &dynamic.Node{ID: "test", TenantID: "test", Name: "test"}
	session1 := NewMockSession()
	builder.addSession(node.ID, session1)

	cmd := &Command{
		Type: UPDATE,
		Node: node,
	}
	manager.executeCommand(context.Background(), cmd)
}

func TestRunSession(t *testing.T) {
	provider := &MockProvider{}
	builder := &MockBuilder{
		mux:      &sync.Mutex{},
		sessions: make(map[string]*MockSession),
	}
	manager := NewManager(builder, provider)

	node := &dynamic.Node{ID: keyError}
	ctx := context.WithValue(context.Background(), keyError, true) // nolint

	go func() {
		manager.runSession(ctx, node)
	}()
}
