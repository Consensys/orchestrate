package stress

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient"
	"github.com/spf13/viper"

	orchestrateclient "github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	utils2 "github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/ConsenSys/orchestrate/tests/service/stress/assets"
	"github.com/ConsenSys/orchestrate/tests/service/stress/units"
	"github.com/ConsenSys/orchestrate/tests/utils/chanregistry"
	"github.com/Shopify/sarama"
)

type WorkLoadTest func(context.Context, *units.WorkloadConfig, orchestrateclient.OrchestrateClient, *chanregistry.ChanRegistry) error

type WorkLoadService struct {
	cfg      *Config
	client   orchestrateclient.OrchestrateClient
	ec       ethclient.MultiClient
	producer sarama.SyncProducer
	chanReg  *chanregistry.ChanRegistry
	items    []*workLoadItem
	cancel   context.CancelFunc
}

type workLoadItem struct {
	iteration int
	threads   int
	name      string
	call      WorkLoadTest
}

// TODO: make it customizable by ENVs
const (
	nAccounts              = 20
	nPrivGroupPerChain     = 10
	waitForEnvelopeTimeout = time.Minute * 2
)

var artifacts = []string{"SimpleToken", "Counter", "ERC20", "ERC721"}

// Init initialize Cucumber service
func NewService(cfg *Config,
	chanReg *chanregistry.ChanRegistry,
	client orchestrateclient.OrchestrateClient,
	ec ethclient.MultiClient,
	producer sarama.SyncProducer,
) *WorkLoadService {
	return &WorkLoadService{
		cfg:      cfg,
		chanReg:  chanReg,
		client:   client,
		ec:       ec,
		producer: producer,
		items: []*workLoadItem{
			{cfg.Iterations, cfg.Concurrency, "BatchDeployContract", units.BatchDeployContractTest},
			{cfg.Iterations, cfg.Concurrency, "BatchPrivateTxsTest", units.BatchPrivateTxsTest},
		},
	}
}

func (c *WorkLoadService) Run(ctx context.Context) error {
	logger := log.FromContext(ctx).WithField("iteration", c.cfg.Iterations).
		WithField("concurrency", c.cfg.Concurrency).
		WithField("timeout", c.cfg.Timeout.String())

	logger.Info("stress test started")

	ctx, c.cancel = context.WithTimeout(ctx, c.cfg.Timeout)

	cctx, err := c.preRun(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var gerr error

	unitCfg := units.NewWorkloadConfig(cctx, waitForEnvelopeTimeout)
	for _, item := range c.items {
		wg.Add(1)
		go func(it *workLoadItem) {
			defer wg.Done()
			err := c.run(cctx, it, unitCfg)
			if err != nil {
				if gerr == nil {
					gerr = err
				}
				c.Stop()
			}
		}(item)
	}

	log.FromContext(ctx).Info("waiting for jobs to complete...")
	wg.Wait()

	return c.postRun(cctx)
}

func (c *WorkLoadService) Stop() {
	c.cancel()
}

func (c *WorkLoadService) preRun(ctx context.Context) (context.Context, error) {
	proxyHost := viper.GetString(orchestrateclient.URLViperKey)
	var err error

	for idx := 0; idx <= nAccounts; idx++ {
		ctx, err = assets.CreateNewAccount(ctx, c.client)
		if err != nil {
			return ctx, err
		}
	}

	for _, contractName := range artifacts {
		ctx, err = assets.RegisterNewContract(ctx, c.client, c.cfg.ArtifactPath, contractName)
		if err != nil {
			return ctx, err
		}
	}

	nBesuNodes := len(c.cfg.gData.Nodes.Besu)
	for idx := 0; idx < nBesuNodes; idx++ {
		besuNode := c.cfg.gData.Nodes.Besu[idx]
		chainName := fmt.Sprintf("besu_%d-%s", idx, utils2.RandString(5))
		var cUUID string
		ctx, cUUID, err = assets.RegisterNewChain(ctx, c.client, c.ec, proxyHost, chainName, &besuNode)
		if err != nil {
			return ctx, err
		}

		privNodeAddress := []string{}
		for jdx := 0; jdx < nBesuNodes; jdx++ {
			besuNode2 := c.cfg.gData.Nodes.Besu[jdx]
			if idx != jdx {
				privNodeAddress = append(privNodeAddress, besuNode2.PrivateAddress...)
			}
		}

		for jdx := 0; jdx < nPrivGroupPerChain; jdx++ {
			ctx, err = assets.CreatePrivateGroup(ctx, c.ec, utils2.GetProxyURL(proxyHost, cUUID), besuNode.PrivateAddress,
				utils2.RandShuffle(privNodeAddress))
			if err != nil {
				return ctx, err
			}
		}
	}

	return ctx, nil
}

func (c *WorkLoadService) postRun(ctx context.Context) error {
	logger := log.FromContext(ctx)
	chains := assets.ContextChains(ctx)

	var gerr error
	for _, chain := range chains {
		err := assets.DeregisterChain(ctx, c.client, &chain)
		if err != nil {
			gerr = errors.CombineErrors(gerr, err)
			logger.WithError(err).Error("failed to remove chain")
		}
	}

	return gerr
}

func (c *WorkLoadService) run(ctx context.Context, test *workLoadItem, cfg *units.WorkloadConfig) error {
	logger := log.FromContext(ctx).WithField("name", test.name)
	logger.Info("test unit started")

	var wg sync.WaitGroup
	started := time.Now()
	buffer := make(chan bool, test.threads)

	var gerr error
	for idx := 0; idx <= test.iteration && gerr == nil; idx++ {
		buffer <- true
		wg.Add(1)
		go func(idx int) {
			err := test.call(ctx, cfg, c.client, c.chanReg)
			if err != nil {
				if gerr == nil {
					gerr = err
				}
				c.Stop()
			}
			wg.Done()
			<-buffer
			if idx != 0 && idx%100 == 0 {
				logger.WithField("iteration", idx).WithField("time", time.Since(started).String()).
					Info("iteration completed")
			}
		}(idx)
	}

	wg.Wait()
	return gerr
}
