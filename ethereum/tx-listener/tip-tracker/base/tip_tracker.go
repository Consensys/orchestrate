package base

import (
	"context"
	"math/big"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

// Tracker is a basic chain tracker that consider a block final if a a certain depth
type Tracker struct {
	ec ethclient.ChainLedgerReader

	chainID *big.Int
	conf    *Config
}

// NewTracker creates a new base tracker
func NewTracker(ec ethclient.ChainLedgerReader, chainID *big.Int, conf *Config) *Tracker {
	log.Infof("Creating new tracker %s", chainID.String())

	return &Tracker{
		ec:      ec,
		chainID: chainID,
		conf:    conf,
	}
}

// ChainID returns ID of the tracked chain
func (t *Tracker) ChainID() *big.Int {
	return big.NewInt(0).Set(t.chainID)
}

// HighestBlock returns highest mined & considered final block on the tracked chain
func (t *Tracker) HighestBlock(ctx context.Context) (int64, error) {
	logCtx := logger.WithLogEntry(
		ctx,
		log.NewEntry(log.StandardLogger()).
			WithFields(log.Fields{
				"chain.id": t.chainID.Text(10),
			}),
	)
	header, err := t.ec.HeaderByNumber(logCtx, t.chainID, nil)
	if err != nil {
		return 0, errors.FromError(err).ExtendComponent(component)
	}

	if header.Number.Uint64() <= t.conf.Depth {
		return 0, nil
	}

	return int64(header.Number.Uint64() - t.conf.Depth), nil
}
