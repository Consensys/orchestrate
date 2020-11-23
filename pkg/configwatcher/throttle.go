package configwatcher

import (
	"context"
	"time"

	"github.com/eapache/channels"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func Throttle(ctx context.Context, throttleDuration time.Duration, in <-chan interface{}, out chan<- interface{}) {
	ring := channels.NewRingChannel(1)
	defer ring.Close()

	utils.InParallel(
		func() {
			// Feeding output loop
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-ring.Out():
					throttle := time.After(throttleDuration)
					select {
					case out <- msg:
					case <-ctx.Done():
						return
					}
					<-throttle
				}
			}
		},
		func() {
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-in:
					ring.In() <- msg
				}
			}
		},
	)
}
