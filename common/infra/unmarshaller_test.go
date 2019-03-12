package infra

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func newProtoMessage() *trace.Trace {
	return &trace.Trace{
		Sender: &common.Account{Id: "abcde"},
	}
}

func TestTracePbUnmarshaller(t *testing.T) {
	u := TracePbUnmarshaller{}
	traces := make([]*trace.Trace, 0)
	rounds := 1000
	wg := &sync.WaitGroup{}
	for i := 1; i < rounds; i++ {
		traces = append(traces, &trace.Trace{})
		wg.Add(1)
		go func(t *trace.Trace) {
			defer wg.Done()
			u.Unmarshal(newProtoMessage(), t)
		}(traces[len(traces)-1])
	}
	wg.Wait()

	for _, tr := range traces {
		if tr.Sender.Id != "abcde" {
			t.Errorf("TracePbUnmarshaller: expected %q but got %q", "abcde", tr.Sender.Id)
		}
	}
}
