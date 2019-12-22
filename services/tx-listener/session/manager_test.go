package session

import (
	"context"
	"sync"
	"testing"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
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

func (p *MockProvider) Run(ctx context.Context, configInput chan<- *dynamic.Message) error {
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
