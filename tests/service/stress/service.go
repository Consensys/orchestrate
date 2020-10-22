package stress

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	identitymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/client"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/stress/units"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/stress/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/utils/chanregistry"
)

type WorkLoadTest func(context.Context, txscheduler.TransactionSchedulerClient, *chanregistry.ChanRegistry) error

type WorkLoadService struct {
	cfg                    *Config
	chainRegistryClient    chainregistry.ChainRegistryClient
	contractRegistryClient registry.ContractRegistryClient
	txSchedulerClient      txscheduler.TransactionSchedulerClient
	identityClient         identitymanager.IdentityManagerClient
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

// Init initialize Cucumber service
func NewService(cfg *Config,
	chanReg *chanregistry.ChanRegistry,
	chainRegistryClient chainregistry.ChainRegistryClient,
	contractRegistryClient registry.ContractRegistryClient,
	txSchedulerClient txscheduler.TransactionSchedulerClient,
	identityClient identitymanager.IdentityManagerClient,
	producer sarama.SyncProducer,
) *WorkLoadService {
	return &WorkLoadService{
		cfg:                    cfg,
		chanReg:                chanReg,
		chainRegistryClient:    chainRegistryClient,
		contractRegistryClient: contractRegistryClient,
		txSchedulerClient:      txSchedulerClient,
		identityClient:         identityClient,
		producer:               producer,
		items: []*workLoadItem{
			{cfg.Iterations, cfg.Concurrency, "BatchDeployContract", units.BatchDeployContractTest},
		},
	}
}

func (c *WorkLoadService) Run(ctx context.Context) error {
	ctx, c.cancel = context.WithCancel(ctx)

	cctx, err := c.preRun(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(c.items))
	var gerr error

	for _, item := range c.items {
		go func(it *workLoadItem) {
			defer wg.Done()
			err := c.run(cctx, it)
			if err != nil {
				gerr = errors.CombineErrors(gerr, err)
			}
		}(item)
	}

	wg.Wait()
	return gerr
}

func (c *WorkLoadService) Stop() {
	c.cancel()
}

func (c *WorkLoadService) preRun(ctx context.Context) (context.Context, error) {
	accounts := []string{}
	for idx := 0; idx <= utils.NAccounts; idx++ {
		acc, err := utils.CreateNewAccount(ctx, c.identityClient)
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
	log.FromContext(ctx).Debugf("Started \"%s\": (%d/%d)", strings.ToUpper(test.name), test.iteration, test.threads)
	var wg sync.WaitGroup
	wg.Add(test.iteration)
	buffer := make(chan bool, test.threads)
	var gerr error
	for idx := 1; idx <= test.iteration && gerr == nil; idx++ {
		buffer <- true
		go func() {
			err := test.call(ctx, c.txSchedulerClient, c.chanReg)
			if err != nil {
				gerr = errors.CombineErrors(gerr, err)
			}
			wg.Done()
			<-buffer
		}()
	}

	wg.Wait()
	return gerr
}
