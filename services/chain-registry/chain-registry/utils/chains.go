package utils

import (
	"context"
	"math/big"

	"github.com/containous/traefik/v2/pkg/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
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

func GetChainTip(ctx context.Context, ec ethclient.ChainLedgerReader, uris []string) (head uint64, err error) {
	var header *types.Header
	for _, uri := range uris {
		header, err = ec.HeaderByNumber(ctx, uri, nil)
		if err == nil {
			return header.Number.Uint64(), nil
		}
		log.FromContext(ctx).WithError(err).Warnf("failed to fetch chain id for URL %s", uri)
	}
	return
}

func GetAddressBalance(ctx context.Context, ec ethclient.ChainStateReader, uris []string, address string) (*big.Int, error) {
	var balance *big.Int
	var err error

	if len(uris) == 0 {
		errMessage := invalidURLListErrMsg
		log.FromContext(ctx).Errorf(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	// All URLs must be valid and we return the head of the latest one
	for _, uri := range uris {
		balance, err = ec.BalanceAt(ctx, uri, ethcommon.HexToAddress(address), nil)
		if err != nil {
			log.FromContext(ctx).
				WithField("address", address).
				WithField("url", uri).
				WithError(err).
				Error("failed to fetch balance")

			continue
		}

		return balance, nil
	}

	errMessage := notReachableURLErrMsg
	log.FromContext(ctx).WithError(err).Errorf(errMessage)
	return nil, errors.ConnectionError(errMessage)
}
