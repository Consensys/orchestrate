package crafter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/tx"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/tx-sender/tx-sender/nonce"
	usecases "github.com/consensys/orchestrate/services/tx-sender/tx-sender/use-cases"
	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const estimationGasError = "cannot estimate gas usage"
const craftTransactionComponent = "use-cases.craft-transaction"

const mediumPriorityString = "1500000000" // 1.5 gwei
const thresholdString = "500000000"       // 0.5 gwei

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
	if string(job.Type) == tx.JobType_ETH_EEA_MARKING_TX.String() {
		if err := uc.craftEEAMarkingTx(ctx, job); err != nil {
			return err
		}
	}

	if job.Transaction.TransactionType == "" {
		if err := uc.craftTransactionType(ctx, job); err != nil {
			return err
		}
	}

	switch job.Transaction.TransactionType {
	case entities.LegacyTxType:
		if job.Transaction.GasPrice == "" {
			if err := uc.craftGasPrice(ctx, job); err != nil {
				return err
			}
		}
	default:
		// We MUST recalculate gasFeeCap for child jobs
		if job.Transaction.GasFeeCap == "" || job.InternalData.ParentJobUUID == job.UUID {
			if err := uc.craftDynamicFeePrice(ctx, job); err != nil {
				return err
			}
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

	if string(job.Type) == tx.JobType_ETH_EEA_PRIVATE_TX.String() {
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
	if string(job.Type) == tx.JobType_ETH_EEA_MARKING_TX.String() {
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

	if string(job.Type) == tx.JobType_ETH_EEA_PRIVATE_TX.String() {
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

	job.Transaction.TransactionType = entities.LegacyTxType
	logger.WithField("value", job.Transaction.GasPrice).Debug("crafted gas price")
	return nil
}

func (uc *craftTxUseCase) craftTransactionType(_ context.Context, job *entities.Job) error {
	switch {
	case job.Transaction.GasPrice != "" || job.InternalData.OneTimeKey:
		job.Transaction.TransactionType = entities.LegacyTxType
	case job.Transaction.GasTipCap != "" || job.Transaction.GasFeeCap != "":
		job.Transaction.TransactionType = entities.DynamicFeeTxType
	}

	return nil
}

func (uc *craftTxUseCase) craftDynamicFeePrice(ctx context.Context, job *entities.Job) error {
	logger := uc.logger.WithContext(ctx)

	if string(job.Type) == tx.JobType_ETH_EEA_PRIVATE_TX.String() {
		logger.Debug("skip gas dynamic fee estimation. EEA private transaction")
		return nil
	}

	proxyURL := utils.GetProxyURL(uc.chainRegistryURL, job.ChainUUID)
	feeHistory, err := uc.ec.FeeHistory(ctx, proxyURL, 1, "latest")
	if err != nil {
		logger.WithError(err).Debug("failed to fetch feeHistory. Fallback to craft GasPrice")
		return uc.craftGasPrice(ctx, job)
	}

	if feeHistory == nil || len(feeHistory.BaseFeePerGas) == 0 {
		logger.Debug("cannot extract base fee. Fallback to craft GasPrice")
		return uc.craftGasPrice(ctx, job)
	}

	nextBlockBaseFeePerGas := feeHistory.BaseFeePerGas[len(feeHistory.BaseFeePerGas)-1].ToInt()
	if nextBlockBaseFeePerGas.String() == "0" {
		logger.Debug("skip gas dynamic fee. Zero base fee is not allowed")
		return uc.craftGasPrice(ctx, job)
	}

	var priorityFee *big.Int
	if job.Transaction.GasTipCap == "" {
		mediumPriority, _ := new(big.Int).SetString(mediumPriorityString, 10) // 1.5 gwei
		threshold, _ := new(big.Int).SetString(thresholdString, 10)           // 0.5 gwei

		switch job.InternalData.Priority {
		case utils.PriorityVeryLow:
			priorityFee = new(big.Int).Sub(mediumPriority, new(big.Int).Mul(threshold, big.NewInt(2))) // 0.5 gwei
		case utils.PriorityLow:
			priorityFee = new(big.Int).Sub(mediumPriority, threshold) // 1 gwei
		case utils.PriorityMedium:
			priorityFee = mediumPriority // 1.5 gwei
		case utils.PriorityHigh:
			priorityFee = new(big.Int).Add(mediumPriority, threshold) // 2 gwei
		case utils.PriorityVeryHigh:
			priorityFee = new(big.Int).Add(mediumPriority, new(big.Int).Mul(threshold, big.NewInt(2))) // 2.5 gwei
		default:
			priorityFee = mediumPriority
		}

		job.Transaction.GasTipCap = priorityFee.String()
	} else {
		priorityFee, _ = new(big.Int).SetString(job.Transaction.GasTipCap, 10)
	}

	job.Transaction.GasFeeCap = new(big.Int).Add(nextBlockBaseFeePerGas, priorityFee).String()
	job.Transaction.TransactionType = entities.DynamicFeeTxType

	logger.WithField("base", nextBlockBaseFeePerGas).WithField("tip", priorityFee).
		Debug("crafted dynamic fees")
	return nil
}
