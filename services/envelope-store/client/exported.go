package client

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/viper"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

var initOnce = &sync.Once{}
var client *Client

func Init(ctx context.Context) {
	cfg := NewConfigFromViper(viper.GetViper()) 
	initOnce.Do(func() {
		var err error
		client, err = NewClient(ctx, &cfg)
		if err != nil {
			log.FromContext(ctx).WithError(err).Fatalf("Could not create envelope store application")
		}
	})
}

func GlobalEnvelopeStoreClient() svc.EnvelopeStoreClient {
	return client.srv
}

// SetGlobalConfig sets Sarama global configuration
func SetGlobalEnvelopeStoreClient(srv svc.EnvelopeStoreClient) {
	client.srv = srv
}