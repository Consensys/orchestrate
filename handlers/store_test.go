package handlers

import (
	"sync"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/core/infra"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core/types"
)

type DummyStore struct {
	mux    *sync.Mutex
	stored []*types.Trace
}

func (s *DummyStore) Store(t *types.Trace) error {
	s.mux.Lock()
	s.stored = append(s.stored, t)
	s.mux.Unlock()
	return nil
}

var testID = "abc"

func newStoreTestMessage() *tracepb.Trace {
	var pb tracepb.Trace
	pb.Sender = &tracepb.Account{Id: testID}
	return &pb
}

func TestStore(t *testing.T) {
	// Create worker
	w := infra.NewWorker(100)

	// Create Sarama loader
	h := Loader(&TraceProtoUnmarshaller{})
	w.Use(h)

	// Register mock handler
	mockH := NewMockHandler(50)
	w.Use(mockH.Handler())

	// Register sender handler
	store := DummyStore{&sync.Mutex{}, []*types.Trace{}}
	w.Use(Store(&store))

	// Create input channel
	in := make(chan interface{})

	// Run worker
	go w.Run(in)

	// Feed input channel and then close it
	rounds := 1000
	for i := 1; i <= rounds; i++ {
		in <- newStoreTestMessage()
	}
	close(in)

	// Wait for worker to be done
	<-w.Done()

	if len(mockH.handled) != rounds {
		t.Errorf("Store: expected %v rounds but got %v", rounds, len(mockH.handled))
	}

	for _, tr := range store.stored {
		if tr.Sender().ID != testID {
			t.Errorf("store: expected ID %v but got %v", testID, tr.Sender().ID)
		}
	}
}
