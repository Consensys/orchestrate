package utils

import (
	"context"
	"math/big"

	"github.com/containous/traefik/v2/pkg/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient"
)

const invalidURLListErrMsg = "invalid URLs list"
const notReachableURLErrMsg = "All URLs in the list are unreachable"

func VerifyURLs(ctx context.Context, ec ethclient.ChainSyncReader, uris []string) error {
	var prevChainID *big.Int

	for i, uri := range uris {
		chainID, err := ec.Network(ctx, uri)
		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("failed to fetch chain id for URL %s", uri)
			return err
		}

		if i > 0 && chainID != prevChainID {
			errMessage := "URLs in the list point to different networks"
			log.FromContext(ctx).Errorf(errMessage)
			return errors.InvalidParameterError(errMessage)
		}

		prevChainID = chainID
	}

	return nil
}

func GetChainID(ctx context.Context, ec ethclient.ChainSyncReader, uris []string) (*big.Int, error) {
	var chainID *big.Int
	var err error

	if len(uris) == 0 {
		errMessage := invalidURLListErrMsg
		log.FromContext(ctx).Errorf(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	// All URLs must be valid and we return the head of the latest one
	for _, uri := range uris {
		chainID, err = ec.Network(ctx, uri)
		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("failed to fetch chain id for URL %s", uri)
			continue
		}

		return chainID, nil
	}

	errMessage := notReachableURLErrMsg
	log.FromContext(ctx).WithError(err).Errorf(errMessage)
	return nil, errors.InvalidParameterError(errMessage)
}

func GetChainTip(ctx context.Context, ec ethclient.ChainLedgerReader, uris []string) uint64 {
	for _, uri := range uris {
		head, err := ec.HeaderByNumber(ctx, uri, nil)
		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("failed to fetch chain id for URL %s", uri)
			continue
		}

		return head.Number.Uint64()
	}

	return 0
}

func GetAddressBalance(ctx context.Context, ec ethclient.ChainStateReader, uris []string, address ethcommon.Address) (*big.Int, error) {
	var balance *big.Int
	var err error

	if len(uris) == 0 {
		errMessage := invalidURLListErrMsg
		log.FromContext(ctx).Errorf(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	// All URLs must be valid and we return the head of the latest one
	for _, uri := range uris {
		balance, err = ec.BalanceAt(ctx, uri, address, nil)
		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("failed to fetch address %s balance for URL %s", address.Hex(), uri)
			continue
		}

		return balance, nil
	}

	errMessage := notReachableURLErrMsg
	log.FromContext(ctx).WithError(err).Errorf(errMessage)
	return nil, errors.ConnectionError(errMessage)
}
