package crafter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ConsenSys/orchestrate/pkg/ethclient"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/types/tx"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/nonce"
	usecases "github.com/ConsenSys/orchestrate/services/tx-sender/tx-sender/use-cases"
	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const estimationGasError = "cannot estimate gas usage"
const craftTransactionComponent = "use-cases.craft-transaction"

type craftTxUseCase struct {
	nonceManager     nonce.Manager
	ec               ethclient.MultiClient
	chainRegistryURL string
	logger           *log.Logger
}

func NewCraftTransactionUseCase(ec ethclient.MultiClient, chainRegistryURL string, nonceManager nonce.Manager) usecases.CraftTransactionUseCase {
	return &craftTxUseCase{
		ec:               ec,
		chainRegistryURL: chainRegistryURL,
		nonceManager:     nonceManager,
		logger:           log.NewLogger().SetComponent(craftTransactionComponent),
	}
}

func (uc *craftTxUseCase) Execute(ctx context.Context, job *entities.Job) error {
	logger := uc.logger.WithContext(ctx)

	if job.InternalData.ParentJobUUID == job.UUID {
		logger.Debug("skip crafting for job resending")
		return nil
	}

	if string(job.Type) == tx.JobType_ETH_ORION_MARKING_TX.String() {
		if err := uc.craftEEAMarkingTx(ctx, job); err != nil {
			return err
		}
	}

	if job.Transaction.GasPrice == "" {
		if err := uc.craftGasPrice(ctx, job); err != nil {
			return err
		}
	}

	if job.Transaction.Gas == "" {
		if err := uc.craftGasEstimation(ctx, job); err != nil {
			return err
		}
	}

	if job.Transaction.Nonce == "" {
		if err := uc.craftNonce(ctx, job); err != nil {
			return err
		}
	}

	return nil
}

func (uc *craftTxUseCase) craftNonce(ctx context.Context, job *entities.Job) error {
	if job.InternalData.OneTimeKey || string(job.Type) == tx.JobType_ETH_TESSERA_PRIVATE_TX.String() {
		job.Transaction.Nonce = "0"
	} else {
		n, err := uc.nonceManager.GetNonce(ctx, job)
		if err != nil {
			return err
		}
		job.Transaction.Nonce = fmt.Sprintf("%d", n)
	}

	uc.logger.WithContext(ctx).WithField("value", job.Transaction.Nonce).Debug("crafted transaction nonce")
	return nil
}

func (uc *craftTxUseCase) craftEEAMarkingTx(ctx context.Context, job *entities.Job) error {
	logger := uc.logger.WithContext(ctx)
	logger.Debug("crafting EEA precompiled contract address")

	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	privPContractAddr, err := uc.ec.EEAPrivPrecompiledContractAddr(ctx, proxyURL)
	if err != nil {
		errMsg := "cannot retrieve EEA precompiled contract address"
		logger.WithError(err).Errorf(errMsg)
		return err
	}

	job.Transaction.To = privPContractAddr.String()
	logger.WithField("value", privPContractAddr.String()).Debug("crafted EEA precompiled contract address to")
	return nil
}

func (uc *craftTxUseCase) craftGasEstimation(ctx context.Context, job *entities.Job) error {
	logger := uc.logger.WithContext(ctx)

	if string(job.Type) == tx.JobType_ETH_ORION_EEA_TX.String() {
		logger.Debug("skip gas estimation for eea private transaction")
		return nil
	}

	call := &ethereum.CallMsg{}
	if job.InternalData.OneTimeKey {
		call.From = ethcommon.HexToAddress("0x1")
	} else {
		call.From = ethcommon.HexToAddress(job.Transaction.From)
	}

	if job.Transaction.To != "" {
		toAddr := ethcommon.HexToAddress(job.Transaction.To)
		call.To = &toAddr
	}

	if job.Transaction.Value != "" {
		call.Value, _ = new(big.Int).SetString(job.Transaction.Value, 10)
	}

	if job.Transaction.Data != "" {
		var err error
		call.Data, err = hexutil.Decode(job.Transaction.Data)
		if err != nil {
			logger.WithError(err).Error(estimationGasError)
			return err
		}
	}

	// We update the data to an arbitrary hash
	// to avoid errors raised on eth_estimateGas on Besu 1.5.4 & 1.5.5
	if string(job.Type) == tx.JobType_ETH_ORION_MARKING_TX.String() {
		call.Data = hexutil.MustDecode("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	}

	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	gasEstimated, err := uc.ec.EstimateGas(ctx, proxyURL, call)
	if err != nil {
		logger.WithError(err).Error(estimationGasError)
		return err
	}

	job.Transaction.Gas = fmt.Sprintf("%d", gasEstimated)
	logger.WithField("value", job.Transaction.Gas).Debug("crafted gas estimation")
	return nil
}

func (uc *craftTxUseCase) craftGasPrice(ctx context.Context, job *entities.Job) error {
	logger := uc.logger.WithContext(ctx)

	if string(job.Type) == tx.JobType_ETH_ORION_EEA_TX.String() {
		logger.Debug("skip gas estimation for eea private transaction")
		return nil
	}

	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	gasPrice, err := uc.ec.SuggestGasPrice(ctx, proxyURL)
	if err != nil {
		logger.WithError(err).Error("cannot suggest gas price")
		return err
	}

	switch job.InternalData.Priority {
	case utils.PriorityVeryLow:
		job.Transaction.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(6)).Div(gasPrice, big.NewInt(10)).String()
	case utils.PriorityLow:
		job.Transaction.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(8)).Div(gasPrice, big.NewInt(10)).String()
	case utils.PriorityMedium:
		job.Transaction.GasPrice = gasPrice.String()
	case utils.PriorityHigh:
		job.Transaction.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(12)).Div(gasPrice, big.NewInt(10)).String()
	case utils.PriorityVeryHigh:
		job.Transaction.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(14)).Div(gasPrice, big.NewInt(10)).String()
	default:
		job.Transaction.GasPrice = gasPrice.String()
	}

	logger.WithField("value", job.Transaction.GasPrice).Debug("crafted gas price")
	return nil
}
