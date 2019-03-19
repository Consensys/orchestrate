package infra

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	store "gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
)

// Infra infrastructure elements of the application
type Infra struct {
	ctx context.Context

	Unmarshaller services.Unmarshaller

	TxSender services.TxSender
	Store    store.TraceStore
	Mec      *ethclient.MultiEthClient

	// TODO: we still have some coupling with Sarama (it should be removed)
	SaramaClient sarama.Client

	closeOnce *sync.Once
	cancel    func()
}

// NewInfra creates a new infrastructure
func NewInfra() *Infra {
	ctx, cancel := context.WithCancel(context.Background())
	i := &Infra{
		ctx:       ctx,
		cancel:    cancel,
		closeOnce: &sync.Once{},
	}

	return i
}

// Init intilialize infrastructure
func (infra *Infra) Init() {
	wait := &sync.WaitGroup{}
	initStore(infra)
	wait.Add(2)
	go initSarama(infra, wait)
	go initEthereum(infra, wait)
	wait.Wait()
}

// Close infra
func (infra *Infra) Close() {
	infra.closeOnce.Do(infra.cancel)
}
