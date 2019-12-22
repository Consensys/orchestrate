package kafka

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
)

type Hook struct {
	Node *dynamic.Node
}

func (hk *Hook) logger(ctx context.Context) logrus.FieldLogger {
	return log.FromContext(ctx).WithFields(logrus.Fields{
		"hook": "kafka",
	})
}

func (hk *Hook) Receipt(ctx context.Context, node *dynamic.Node, block *ethtypes.Block, receipt *ethtypes.Receipt) error {
	hk.logger(ctx).WithFields(logrus.Fields{
		"receipt.txHash": receipt.TxHash.Hex(),
	}).Debugf("Receipt processed")
	return nil
}

func (hk *Hook) Block(ctx context.Context, node *dynamic.Node, block *ethtypes.Block) error {
	hk.logger(ctx).WithFields(logrus.Fields{
		"block.Hash": block.Hash().Hex(),
	}).Debugf("Block processed")
	return nil
}
