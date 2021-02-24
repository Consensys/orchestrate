// +build unit

package dispatcher

import (
	"fmt"
	"testing"

	"github.com/ConsenSys/orchestrate/tests/service/e2e/cucumber/alias"

	"github.com/Shopify/sarama"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/stretchr/testify/assert"
	broker "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	"github.com/ConsenSys/orchestrate/pkg/engine"
	"github.com/ConsenSys/orchestrate/pkg/types/tx"
	"github.com/ConsenSys/orchestrate/tests/utils/chanregistry"
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
	txctx.Logger = log.NewLogger()
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

	txctx.Envelope = txctx.Envelope.SetJobUUID("jobUUID")

	return txctx
}

func TestDispatcher(t *testing.T) {
	reg := chanregistry.NewChanRegistry()
	ch := make(chan *tx.Envelope, 10)
	reg.Register("known-key", ch)
	reg.Register("tx.decoded/"+alias.ExternalTxLabel, ch)
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

	// external tx
	txctx = makeContext("known-key", "")
	txctx.Envelope = txctx.Envelope.SetJobUUID("")
	h(txctx)
	select {
	case e := <-ch:
		assert.Equal(t, txctx.Envelope, e, "#1: Envelope should match")
	default:
		t.Errorf("#1: Envelope should have been dispatched")
	}
}
