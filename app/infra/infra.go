package infra

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
)

// Infra infrastructure elements of the application
type Infra struct {
	ctx context.Context

	Unmarshaller services.Unmarshaller
	Producer     services.Producer

	Signer services.TxSigner

	Mec *ethclient.MultiEthClient

	// TODO: we still have some coupling with Sarama (it should be removed)
	SaramaClient   sarama.Client
	SaramaProducer sarama.SyncProducer

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
	wait.Add(1)
	go initSarama(infra, wait)
	wait.Wait()
	initSigner(infra)
}

// Close infra
func (infra *Infra) Close() {
	infra.closeOnce.Do(infra.cancel)
}
