package sarama

import (
	"context"

	"github.com/Shopify/sarama"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Pipe take a channel of sarama.ConsumerMessage and pipes it into a channel of engine.Msg
//
// Pipe will stop forwarding messages either
// - sarama channel is closed
// - ctx has been canceled
func Pipe(ctx context.Context, saramaChan <-chan *sarama.ConsumerMessage) <-chan engine.Msg {
	msgChan := make(chan engine.Msg)

	// Start a goroutine that pipe messages
	go func() {
	pipeLoop:
		for {
			select {
			case msg, ok := <-saramaChan:
				if !ok {
					// Sarama channel has been closed so we exit loop
					break pipeLoop
				}
				msgChan <- &Msg{*msg}
			case <-ctx.Done():
				// Context has been cancel so we exit loop
				break pipeLoop
			}
		}
		close(msgChan)
	}()

	return msgChan
}
