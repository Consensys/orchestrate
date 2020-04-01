// +build unit

package dispatcher

import (
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/tx"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/chanregistry"
)

func testKeyOf1(txctx *engine.TxContext) (string, error) {
	key, ok := txctx.Get("key1").(string)
	if !ok {
		return "", fmt.Errorf("unknown key")
	}
	return key, nil
}

func testKeyOf2(txctx *engine.TxContext) (string, error) {
	key, ok := txctx.Get("key2").(string)
	if !ok {
		return "", fmt.Errorf("unknown key")
	}
	return key, nil
}

func makeContext(key1, key2 string) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	txctx.In = &broker.Msg{
		ConsumerMessage: sarama.ConsumerMessage{
			Topic: "testTopic",
		},
	}

	if key1 != "" {
		txctx.Set("key1", key1)
	}

	if key2 != "" {
		txctx.Set("key2", key2)
	}

	return txctx
}

func TestDispatcher(t *testing.T) {
	reg := chanregistry.NewChanRegistry()
	ch := make(chan *tx.Envelope, 10)
	reg.Register("known-key", ch)
	h := Dispatcher(reg, testKeyOf1, testKeyOf2)

	// Handle context
	txctx := makeContext("known-key", "")
	h(txctx)
	select {
	case e := <-ch:
		assert.Equal(t, txctx.Envelope, e, "#1: Envelope should match")
	default:
		t.Errorf("#1: Envelope should have been dispatched")
	}

	// Handle context
	txctx = makeContext("", "known-key")
	h(txctx)
	select {
	case e := <-ch:
		assert.Equal(t, txctx.Envelope, e, "#2: Envelope should match")
	default:
		t.Errorf("#2: Envelope should have been dispatched")
	}

	// Handle context
	txctx = makeContext("unknown-key", "")
	h(txctx)

	select {
	case <-ch:
		t.Errorf("#3: No envelope should have been dispatched")
	default:
	}
}
