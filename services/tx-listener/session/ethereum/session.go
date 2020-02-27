package ethereum

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	ethclientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session"
	hook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/hooks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset"
)

type fetchedBlock struct {
	block    *ethtypes.Block
	receipts []*ethtypes.Receipt
}

type EthClient interface {
	ethclient.ChainLedgerReader
	ethclient.ChainSyncReader
}

type Session struct {
	Chain *dynamic.Chain

	ec EthClient

	hook    hook.Hook
	offsets offset.Manager

	// Listening session
	trigger         chan struct{}
	blockPosition   uint64
	currentChainTip uint64

	// Channel stacking blocks waiting for receipts to be fetched
	fetchedBlocks chan *Future

	errors chan error
}

func (s *Session) Run(ctx context.Context) error {
	return s.run(ctx)
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

	select {
	// Wait for an error or for context to be canceled
	case err = <-s.errors:
		if err != context.DeadlineExceeded && err != context.Canceled {
			log.FromContext(ctx).WithError(err).Errorf("error while listening")
			// If we get an error we cancel execution
			cancel()
		}
	case <-ctx.Done():
	}

	// Drain errors
	go func() {
		for range s.errors {
		}
	}()

	// Wait for goroutines to complete and close session
	wg.Wait()
	s.close()

	return
}

func (s *Session) init(ctx context.Context) error {
	err := s.initChainID(ctx)
	if err != nil {
		return err
	}

	err = s.initPosition(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) initChainID(ctx context.Context) error {
	chain, err := s.ec.Network(ctx, s.Chain.URL)
	if err != nil {
		return err
	}
	s.Chain.ChainID = chain
	return nil
}

func (s *Session) initPosition(ctx context.Context) error {
	blockPosition, err := s.offsets.GetLastBlockNumber(ctx, s.Chain)
	if err != nil {
		return err
	}

	// if blockPosition and startingBlock are different then we have already started listening to that chain
	// and we start at next block since the current one is already processed
	if blockPosition != s.Chain.Listener.StartingBlock {
		blockPosition++
	}

	s.blockPosition = blockPosition

	return nil
}

func (s *Session) listen(ctx context.Context) {
	log.FromContext(ctx).WithField("block.start", s.blockPosition).Infof("start listening")
	ticker := time.NewTicker(s.Chain.Listener.Backoff)
listeningLoop:
	for {
		select {
		case <-ctx.Done():
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
					break listeningLoop
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

	log.FromContext(ctx).WithField("block.stop", s.blockPosition).Infof("Stopped listening")
}

func (s *Session) callHooks(ctx context.Context) {
	var err error
	for futureBlock := range s.fetchedBlocks {
		select {
		case res := <-futureBlock.Result():
			err = s.callHook(ctx, res.(*fetchedBlock))
		case e := <-futureBlock.Err():
			if err == nil {
				err = e
			}
		}

		// Close future
		futureBlock.Close()

		if err != nil {
			s.errors <- err
		}
	}
}

func (s *Session) callHook(ctx context.Context, block *fetchedBlock) error {
	// Call hook
	err := s.hook.AfterNewBlock(ctx, s.Chain, block.block, block.receipts)
	if err == nil {
		// Update last block processed
		err = s.offsets.SetLastBlockNumber(ctx, s.Chain, block.block.NumberU64())
	}
	return err
}

func (s *Session) fetchBlock(ctx context.Context, blockPosition uint64) *Future {
	return NewFuture(func() (interface{}, error) {
		blck, err := s.ec.BlockByNumber(
			ethclientutils.RetryNotFoundError(ctx, true),
			s.Chain.URL,
			big.NewInt(int64(blockPosition)),
		)
		if err != nil {
			log.FromContext(ctx).WithError(err).WithField("block.number", blockPosition).Errorf("failed to fetch block")
			return nil, err
		}

		block := &fetchedBlock{block: blck}

		// Fetch receipt for every transactions
		futureReceipts := []*Future{}
		for _, tx := range blck.Transactions() {
			futureReceipts = append(futureReceipts, s.fetchReceipt(ctx, tx.Hash()))
		}

		for _, futureReceipt := range futureReceipts {
			select {
			case e := <-futureReceipt.Err():
				if err == nil {
					err = e
				}
			case res := <-futureReceipt.Result():
				block.receipts = append(block.receipts, res.(*ethtypes.Receipt))
			}

			// Close future
			futureReceipt.Close()
		}

		if err != nil {
			return nil, err
		}

		return block, nil
	})
}

func (s *Session) fetchReceipt(ctx context.Context, txHash ethcommon.Hash) *Future {
	return NewFuture(func() (interface{}, error) {
		receipt, err := s.ec.TransactionReceipt(
			ethclientutils.RetryNotFoundError(ctx, true),
			s.Chain.URL,
			txHash,
		)
		if err != nil {
			log.FromContext(ctx).WithError(err).WithField("tx.hash", txHash.Hex()).Errorf("failed to fetch receipt")
			return nil, err
		}
		return receipt, nil
	})
}

func (s *Session) getChainTip(ctx context.Context) (tip uint64, err error) {
	head, err := s.ec.HeaderByNumber(
		ethclientutils.RetryNotFoundError(ctx, true),
		s.Chain.URL,
		nil,
	)
	if err != nil {
		log.FromContext(ctx).WithError(err).Errorf("failed to fetch chain head")
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

func (s *Session) close() {
	close(s.errors)
	close(s.trigger)
}

type SessionBuilder struct {
	hook    hook.Hook
	offsets offset.Manager

	ec EthClient
}

func NewSessionBuilder(hk hook.Hook, offsets offset.Manager, ec EthClient) *SessionBuilder {
	return &SessionBuilder{
		hook:    hk,
		offsets: offsets,
		ec:      ec,
	}
}

func (b *SessionBuilder) NewSession(chain *dynamic.Chain) (session.Session, error) {
	return b.newSession(chain), nil
}

func (b *SessionBuilder) newSession(chain *dynamic.Chain) *Session {
	return &Session{
		Chain:         chain,
		ec:            b.ec,
		hook:          b.hook,
		offsets:       b.offsets,
		trigger:       make(chan struct{}, 1),
		fetchedBlocks: make(chan *Future, 20),
		errors:        make(chan error, 1),
	}
}
