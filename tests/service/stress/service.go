package stress

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/stress/units"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/stress/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils/chanregistry"
)

type WorkLoadTest func(context.Context, *units.WorkloadConfig, orchestrateclient.OrchestrateClient, *chanregistry.ChanRegistry) error

type WorkLoadService struct {
	cfg                    *Config
	chainRegistryClient    chainregistry.ChainRegistryClient
	contractRegistryClient registry.ContractRegistryClient
	client                 orchestrateclient.OrchestrateClient
	producer               sarama.SyncProducer
	chanReg                *chanregistry.ChanRegistry
	items                  []*workLoadItem
	cancel                 context.CancelFunc
}

type workLoadItem struct {
	iteration int
	threads   int
	name      string
	call      WorkLoadTest
}

const (
	nAccounts              = 20
	waitForEnvelopeTimeout = time.Minute * 2 // TODO: make it customizable by ENVs
)

// Init initialize Cucumber service
func NewService(cfg *Config,
	chanReg *chanregistry.ChanRegistry,
	chainRegistryClient chainregistry.ChainRegistryClient,
	contractRegistryClient registry.ContractRegistryClient,
	client orchestrateclient.OrchestrateClient,
	producer sarama.SyncProducer,
) *WorkLoadService {
	return &WorkLoadService{
		cfg:                    cfg,
		chanReg:                chanReg,
		chainRegistryClient:    chainRegistryClient,
		contractRegistryClient: contractRegistryClient,
		client:                 client,
		producer:               producer,
		items: []*workLoadItem{
			{cfg.Iterations, cfg.Concurrency, "BatchDeployContract", units.BatchDeployContractTest},
		},
	}
}

func (c *WorkLoadService) Run(ctx context.Context) error {
	log.FromContext(ctx).WithField("iteration", c.cfg.Iterations).
		WithField("concurrency", c.cfg.Concurrency).
		WithField("timeout", c.cfg.Timeout.String()).
		Info("Stress test started")

	ctx, c.cancel = context.WithTimeout(ctx, c.cfg.Timeout)

	cctx, err := c.preRun(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var gerr error

	for _, item := range c.items {
		if gerr != nil {
			break
		}

		wg.Add(1)
		go func(it *workLoadItem) {
			defer wg.Done()
			err := c.run(cctx, it)
			if err != nil {
				gerr = errors.CombineErrors(gerr, err)
			}
		}(item)
	}

	log.FromContext(ctx).Info("waiting for jobs to complete...")
	wg.Wait()
	return gerr
}

func (c *WorkLoadService) Stop() {
	c.cancel()
}

func (c *WorkLoadService) preRun(ctx context.Context) (context.Context, error) {
	var accounts []string
	for idx := 0; idx <= nAccounts; idx++ {
		acc, err := utils.CreateNewAccount(ctx, c.client)
		if err != nil {
			return ctx, err
		}
		accounts = append(accounts, acc)
	}

	ctx = utils.ContextWithAccounts(ctx, accounts)

	err := utils.RegisterNewContract(ctx, c.contractRegistryClient, c.cfg.ArtifactPath, "SimpleToken")
	if err != nil {
		return ctx, err
	}

	chainName := fmt.Sprintf("besu-%s", utils2.RandomString(5))
	chain, err := utils.RegisterNewChain(ctx, c.chainRegistryClient, chainName, c.cfg.gData.Nodes.BesuOne.URLs)
	if err != nil {
		return ctx, err
	}

	ctx = utils.ContextWithChains(ctx, map[string]*models.Chain{"besu": chain})
	return ctx, nil
}

func (c *WorkLoadService) run(ctx context.Context, test *workLoadItem) error {
	log.FromContext(ctx).Debugf("Started \"%s\"", strings.ToUpper(test.name))
	var wg sync.WaitGroup
	wg.Add(test.iteration)
	started := time.Now()
	buffer := make(chan bool, test.threads)
	unitCfg := units.NewWorkloadConfig(nAccounts, waitForEnvelopeTimeout)

	var gerr error
	for idx := 1; idx <= test.iteration && gerr == nil; idx++ {
		buffer <- true
		go func(idx int) {
			err := test.call(ctx, unitCfg, c.client, c.chanReg)
			if err != nil {
				gerr = errors.CombineErrors(gerr, err)
			}
			wg.Done()
			<-buffer
			if idx%100 == 0 {
				log.FromContext(ctx).Infof("iteration %d completed. Time %s", idx, time.Since(started).String())
			}
		}(idx)
	}

	if gerr != nil {
		return gerr
	}

	wg.Wait()
	return nil
}
