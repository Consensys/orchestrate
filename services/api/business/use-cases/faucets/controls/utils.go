package controls

import (
	"context"
	"math/big"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

func getAddressBalance(ctx context.Context, chainStateReader ethclient.ChainStateReader, uris []string, address ethcommon.Address) (*big.Int, error) {
	for _, uri := range uris {
		balance, err := chainStateReader.BalanceAt(ctx, uri, address, nil)
		if err != nil {
			log.WithContext(ctx).WithField("url", uri).WithError(err).Error("failed to fetch balance")
			continue
		}

		return balance, nil
	}

	return nil, errors.EthConnectionError("all URLs in the list are unreachable")
}
