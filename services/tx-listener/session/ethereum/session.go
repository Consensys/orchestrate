package ethereum

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/tx-listener/metrics"

	"github.com/consensys/orchestrate/services/tx-listener/session"

	"github.com/cenkalti/backoff/v4"
	"github.com/consensys/orchestrate/pkg/errors"
	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/services/tx-listener/dynamic"
	hook "github.com/consensys/orchestrate/services/tx-listener/session/ethereum/hooks"
	"github.com/consensys/orchestrate/services/tx-listener/session/ethereum/offset"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

const MaxTxHashesLength = 30
const component = "tx-listener.session.ethereum"

type Session struct {
	Chain         *dynamic.Chain
	ec            EthClient
	client        orchestrateclient.OrchestrateClient
	hook          hook.Hook
	offsets       offset.Manager
	bckOff        backoff.BackOff
	metrics       metrics.ListenerMetrics
	metricsLabels []string
	// Listening session
	trigger                        chan struct{}
	blockPosition                  uint64
	eeaPrivPrecompiledContractAddr string
	currentChainTip                uint64
	// Channel stacking blocks waiting for receipts to be fetched
	fetchedBlocks chan *Future
	errors        chan error
	logger        *log.Logger
}

func NewSession(
	chain *dynamic.Chain,
	ec EthClient,
	client orchestrateclient.OrchestrateClient,
	callHook hook.Hook,
	offsets offset.Manager,
	m metrics.ListenerMetrics,
) *Session {
	return &Session{
		Chain:   chain,
		ec:      ec,
		client:  client,
		hook:    callHook,
		offsets: offsets,
		bckOff:  backoff.NewConstantBackOff(2 * time.Second),
		metrics: m,
		metricsLabels: []string{
			"chain_uuid", chain.UUID,
		},
		logger: log.NewLogger().SetComponent(component).WithField("chain", chain.UUID),
	}
}

type SessionBuilder struct {
	hook    hook.Hook
	offsets offset.Manager
	ec      EthClient
	client  orchestrateclient.OrchestrateClient
	metrics metrics.ListenerMetrics
}

func NewSessionBuilder(
	hk hook.Hook,
	offsets offset.Manager,
	ec EthClient,
	client orchestrateclient.OrchestrateClient,
	m metrics.ListenerMetrics,
) *SessionBuilder {
	return &SessionBuilder{
		hook:    hk,
		offsets: offsets,
		ec:      ec,
		client:  client,
		metrics: m,
	}
}

func (b *SessionBuilder) NewSession(chain *dynamic.Chain) (session.Session, error) {
	return NewSession(chain, b.ec, b.client, b.hook, b.offsets, b.metrics), nil
}

type fetchedBlock struct {
	block *ethtypes.Block
	jobs  []*entities.Job
}

func (s *Session) Run(ctx context.Context) error {
	ctx = log.With(ctx, s.logger)
	err := backoff.RetryNotify(
		func() error {
			err := s.run(ctx)
			if err == context.DeadlineExceeded || err == context.Canceled || ctx.Err() != nil {
				if err == nil {
					err = ctx.Err()
				}

				s.logger.Debug("exiting listener session...")
				return backoff.Permanent(err)
			}

			return err
		},
		s.bckOff,
		func(err error, duration time.Duration) {
			s.logger.WithError(err).Warnf("error in session listener, rebooting in %v...", duration)
		},
	)

	s.logger.WithError(err).Info("listener session exited")
	return err
}

func (s *Session) run(ctx context.Context) (err error) {
	// Initialize session
	err = s.init(ctx)
	if err != nil {
		return err
	}

	cancelableCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start go-routines
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		s.trig()
		s.listen(cancelableCtx)
		wg.Done()
	}()
	go func() {
		s.callHooks(cancelableCtx)
		wg.Done()
	}()

	// Wait for an error or for context to be canceled
	select {
	case err = <-s.errors:
		cancel()
		break
	case <-ctx.Done():
		cancel()
		break
	}

	// We must drain channels before starting a new session
	go func() {
		for e := range s.errors {
			s.logger.WithError(e).Error("error while listening")
		}
	}()

	s.logger.Debug("waiting for go routines to complete....")
	wg.Wait()

	s.close(ctx)
	return err
}

func (s *Session) init(ctx context.Context) error {
	s.logger.Debug("initializing session listener...")

	err := s.initPosition(ctx)
	if err != nil {
		return err
	}

	s.initEEAPrivPrecompiledContractAddr(ctx)

	s.trigger = make(chan struct{}, 1)
	s.errors = make(chan error, 1)
	s.fetchedBlocks = make(chan *Future, 20)

	return nil
}

func (s *Session) initEEAPrivPrecompiledContractAddr(ctx context.Context) {
	// We ignore the error as this only compiles on networks with { eea:1.0, priv:1.0}
	addr, err := s.ec.EEAPrivPrecompiledContractAddr(ctx, s.Chain.URL)
	if err == nil {
		s.eeaPrivPrecompiledContractAddr = addr.String()
	}
}

func (s *Session) initPosition(ctx context.Context) error {
	blockPosition, err := s.offsets.GetLastBlockNumber(ctx, s.Chain)
	if err != nil {
		return err
	}

	s.metrics.BlockCounter().With(s.metricsLabels...).Add(float64(blockPosition))

	// if blockPosition and startingBlock are different then we have already started listening to that chain
	// and we start at next block since the current one is already processed
	if blockPosition != s.Chain.Listener.StartingBlock {
		blockPosition++
	}

	s.blockPosition = blockPosition

	return nil
}

func (s *Session) listen(ctx context.Context) {
	s.logger.WithField("block_start", s.blockPosition).Info("starting fetch block listener")

	ticker := time.NewTicker(s.Chain.Listener.Backoff)
listeningLoop:
	for {
		select {
		case <-ctx.Done():
			s.logger.WithField("block_stop", s.blockPosition).
				Debug("stopping fetch block listener")
			break listeningLoop
		case <-s.trigger:
			if (s.currentChainTip > 0) && s.blockPosition <= s.currentChainTip {
				s.fetchedBlocks <- s.fetchBlock(ctx, s.blockPosition)
				s.blockPosition++
				s.trig()
			} else {
				//  We are ahead of chain head so we update chain tip
				tip, err := s.getChainTip(ctx)
				if err != nil {
					s.errors <- err
				} else if tip > s.currentChainTip {
					s.currentChainTip = tip
					s.trig()
				}
			}
		case <-ticker.C:
			s.trig()
		}
	}

	// Close channels
	ticker.Stop()
	close(s.fetchedBlocks)

	s.logger.WithField("block_stop", s.blockPosition).
		Info("fetch block listener has been stopped")
}

func (s *Session) callHooks(ctx context.Context) {
	var err error

	for futureBlock := range s.fetchedBlocks {
		select {
		case res := <-futureBlock.Result():
			// We MUST drain array chan and ignore blocks after an error happened
			if err != nil {
				s.logger.
					WithField("blockNumber", res.(*fetchedBlock).block.NumberU64()).
					Warn("ignoring fetched block")
				continue
			}
			err = s.callHook(ctx, res.(*fetchedBlock))
		case e := <-futureBlock.Err():
			if err == nil && e != nil {
				err = e
			}
		}

		// Close future
		futureBlock.Close()

		if err != nil {
			s.errors <- err
		} else {
			s.metrics.BlockCounter().With(s.metricsLabels...).Add(1)
		}
	}

	s.logger.Debug("call hooks loop has been stopped")
}

func (s *Session) callHook(ctx context.Context, block *fetchedBlock) error {
	err := s.hook.AfterNewBlock(ctx, s.Chain, block.block, block.jobs)
	if err != nil {
		return err
	}

	if block.block.NumberU64()%3 == 0 {
		return s.offsets.SetLastBlockNumber(ctx, s.Chain, block.block.NumberU64())
	}

	return nil
}

func (s *Session) fetchBlock(ctx context.Context, blockPosition uint64) *Future {
	return NewFuture(func() (interface{}, error) {
		blck, err := s.ec.BlockByNumber(
			ctx,
			s.Chain.URL,
			big.NewInt(int64(blockPosition)),
		)
		if err != nil {
			errMessage := "failed to fetch block"
			if !errors.IsNotFoundError(err) {
				s.logger.WithError(err).WithField("block_number", blockPosition).Error(errMessage)
			}
			return nil, errors.ConnectionError(errMessage)
		}

		block := &fetchedBlock{block: blck}

		for _, tx := range blck.Transactions() {
			s.logger.WithField("tx_hash", tx.Hash().String()).
				WithField("block_number", blck.NumberU64()).Debug("found transaction in block")
		}

		jobMap, err := s.fetchJobs(ctx, blck.Transactions())
		if err != nil {
			return nil, err
		}

		// TODO: pass batch variable by environment variable
		batch := 20
		for i := 0; i < blck.Transactions().Len(); i += batch {
			j := i + batch
			if j > blck.Transactions().Len() {
				j = blck.Transactions().Len()
			}
			jobs, err := awaitReceipts(s.fetchReceipts(ctx, blck.Transactions()[i:j], jobMap))
			if err != nil {
				return nil, err
			}
			block.jobs = append(block.jobs, jobs...)
		}

		return block, nil
	})
}

func (s *Session) fetchJobs(ctx context.Context, transactions ethtypes.Transactions) (map[string]*entities.Job, error) {
	jobMap := make(map[string]*entities.Job)

	if len(transactions) == 0 {
		return jobMap, nil
	}

	for i := 0; i < transactions.Len(); i += MaxTxHashesLength {
		size := i + MaxTxHashesLength
		if size > transactions.Len() {
			size = transactions.Len()
		}
		currTransactions := transactions[i:size]
		var txHashes []string
		for _, t := range currTransactions {
			txHashes = append(txHashes, t.Hash().String())
		}

		// By design, we will receive 0 or 1 job per tx_hash in the filter because we filter by status PENDING
		jobResponses, err := s.client.SearchJob(ctx, &entities.JobFilters{
			TxHashes:  txHashes,
			ChainUUID: s.Chain.UUID,
			Status:    entities.StatusPending,
		})
		if err != nil {
			s.logger.WithError(err).Error("failed to search jobs")
			return nil, err
		}

		for _, jobResponse := range jobResponses {
			s.logger.WithField("tx_hash", jobResponse.Transaction.Hash).
				WithField("job", jobResponse.UUID).Debug("transaction was matched to a job")

			// Filter by the jobs belonging to same session CHAIN_UUID
			jobMap[jobResponse.Transaction.Hash.String()] = &entities.Job{
				UUID:         jobResponse.UUID,
				ChainUUID:    jobResponse.ChainUUID,
				ScheduleUUID: jobResponse.ScheduleUUID,
				TenantID:     jobResponse.TenantID,
				OwnerID:      jobResponse.OwnerID,
				Type:         jobResponse.Type,
				Labels:       jobResponse.Labels,
				Transaction:  &jobResponse.Transaction,
				CreatedAt:    jobResponse.CreatedAt,
			}
		}
	}

	return jobMap, nil
}

func (s *Session) fetchReceipts(ctx context.Context, transactions ethtypes.Transactions, jobMap map[string]*entities.Job) []*Future {
	var futureJobs []*Future

	for _, blckTx := range transactions {
		switch {
		case isEEAPrivTx(blckTx, s.eeaPrivPrecompiledContractAddr) && isInternalTx(jobMap, blckTx):
			futureJobs = append(futureJobs, s.fetchPrivateReceipt(ctx, jobMap[blckTx.Hash().String()], blckTx.Hash()))
			continue
		case isInternalTx(jobMap, blckTx):
			futureJobs = append(futureJobs, s.fetchReceipt(ctx, jobMap[blckTx.Hash().String()], blckTx.Hash()))
			continue
		case isEEAPrivTx(blckTx, s.eeaPrivPrecompiledContractAddr) && s.Chain.Listener.IsExternalTxEnabled():
			job := &entities.Job{ChainUUID: s.Chain.UUID, Transaction: &entities.ETHTransaction{Hash: utils.ToPtr(blckTx.Hash()).(*ethcommon.Hash)}}
			futureJobs = append(futureJobs, s.fetchPrivateReceipt(ctx, job, blckTx.Hash()))
			continue
		case s.Chain.Listener.IsExternalTxEnabled():
			job := &entities.Job{ChainUUID: s.Chain.UUID, Transaction: &entities.ETHTransaction{Hash: utils.ToPtr(blckTx.Hash()).(*ethcommon.Hash)}}
			futureJobs = append(futureJobs, s.fetchReceipt(ctx, job, blckTx.Hash()))
			continue
		default:
			continue
		}
	}

	return futureJobs
}

func awaitReceipts(futureJobs []*Future) (jobs []*entities.Job, err error) {
	for _, futureJob := range futureJobs {
		select {
		case e := <-futureJob.Err():
			if err == nil {
				err = e
			}
		case res := <-futureJob.Result():
			jobs = append(jobs, res.(*entities.Job))
		}

		// Close future
		futureJob.Close()
	}
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func isInternalTx(jobMap map[string]*entities.Job, transaction *ethtypes.Transaction) bool {
	_, ok := jobMap[transaction.Hash().String()]
	return ok
}

func isEEAPrivTx(transaction *ethtypes.Transaction, eeaPrivPrecompiledContractAddr string) bool {
	if eeaPrivPrecompiledContractAddr == "" {
		return false
	}
	// A enclavekey tx has as To address the pre-deployed smart-contract
	return transaction.To() != nil && transaction.To().String() == eeaPrivPrecompiledContractAddr
}

func (s *Session) fetchReceipt(ctx context.Context, job *entities.Job, txHash ethcommon.Hash) *Future {
	return NewFuture(func() (interface{}, error) {
		logger := s.logger.WithField("tx_hash", txHash.Hex()).WithField("chain", s.Chain.UUID)
		logger.Debug("fetching fetch receipt...")

		receipt, err := s.ec.TransactionReceipt(ctx, s.Chain.URL, txHash)
		if err != nil {
			logger.WithError(err).Error("failed to fetch receipt")
			return nil, err
		}

		// Attach receipt to envelope
		job.Receipt = receipt.
			SetBlockHash(ethcommon.HexToHash(receipt.GetBlockHash())).
			SetBlockNumber(receipt.GetBlockNumber()).
			SetTxIndex(receipt.TxIndex)

		return job, nil
	})
}

func (s *Session) fetchPrivateReceipt(ctx context.Context, job *entities.Job, txHash ethcommon.Hash) *Future {
	return NewFuture(func() (interface{}, error) {
		logger := s.logger.WithField("tx_hash", txHash.Hex()).WithField("chain", s.Chain.UUID)

		logger.Debug("fetching private receipt")

		receipt, err := s.ec.PrivateTransactionReceipt(
			ctx,
			s.Chain.URL,
			txHash,
		)

		// We exit ONLY when we failed to fetch the marking tx receipt, otherwise
		// error is being appended to the envelope
		if err != nil && receipt == nil {
			logger.Error("failed to fetch private receipt")
			return nil, err
		} else if receipt == nil {
			logger.Debug("fetched an empty receipt")
			return nil, nil
		}

		logger.WithField("status", receipt.Status).Debug("private receipt was fetched")

		// Bind the hybrid receipt to the envelope
		job.Receipt = receipt.
			SetBlockHash(ethcommon.HexToHash(receipt.GetBlockHash())).
			SetBlockNumber(receipt.GetBlockNumber()).
			SetTxHash(txHash).
			SetTxIndex(receipt.TxIndex)

		return job, nil
	})
}

func (s *Session) getChainTip(ctx context.Context) (tip uint64, err error) {
	head, err := s.ec.HeaderByNumber(ctx, s.Chain.URL, nil)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch chain head")
		return 0, err
	}

	if head.Number.Uint64() > s.Chain.Listener.Depth {
		tip = head.Number.Uint64() - s.Chain.Listener.Depth
	}

	return
}

func (s *Session) trig() {
	select {
	case s.trigger <- struct{}{}:
	default:
		// already triggered
	}
}

func (s *Session) close(_ context.Context) {
	s.logger.Debug("closing session...")
	close(s.errors)
	close(s.trigger)
}
