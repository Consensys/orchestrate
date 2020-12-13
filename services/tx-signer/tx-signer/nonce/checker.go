package nonce

import (
	"context"
	"crypto/md5"
	"fmt"
	"sort"
	"strconv"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/tx"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/nonce"
)

const fetchNonceErr = "cannot retrieve fetch nonce from chain"

//go:generate mockgen -source=checker.go -destination=mocks/checker.go -package=mocks

type Checker interface {
	Check(ctx context.Context, job *entities.Job) error
	OnFailure(ctx context.Context, job *entities.Job, jobErr error) error
	OnSuccess(ctx context.Context, job *entities.Job) error
}

type checker struct {
	nonce            nonce.Sender
	ethClient        ethclient.MultiClient
	recovery         *RecoveryTracker
	maxRecovery      uint64
	chainRegistryURL string
}

func NewNonceChecker(ec ethclient.MultiClient, nm nonce.Sender, tracker *RecoveryTracker, chainRegistryURL string, maxRecovery uint64) Checker {
	return &checker{
		nonce:            nm,
		ethClient:        ec,
		recovery:         tracker,
		maxRecovery:      maxRecovery,
		chainRegistryURL: chainRegistryURL,
	}
}

func (nc *checker) Check(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	if reason, skip := nc.shouldSkip(job); skip {
		logger.Debugf("skip nonce check. %s", reason)
		return nil
	}

	logger.Debug("checking job nonce")

	nonceKey := partitionKey(job)
	if nonceKey == "" {
		logger.Debug("empty nonceKey, skip nonce check")
		return nil
	}

	// Retrieve last sent nonce from nonce manager
	lastSent, ok, err := nc.nonce.GetLastSent(nonceKey)
	if err != nil {
		errMsg := "cannot retrieve lastSent nonce"
		logger.WithError(err).Error(errMsg)
		return err
	}

	var expectedNonce uint64
	if ok {
		expectedNonce = lastSent + 1
	} else {
		expectedNonce, err = nc.fetchNonceFromChain(ctx, job)
		if err != nil {
			logger.WithError(err).Error(fetchNonceErr)
			return err
		}

		logger.WithField("nonce.pending", expectedNonce).Debug("calibrating nonce")
	}

	txNonce, _ := strconv.ParseUint(job.Transaction.Nonce, 10, 32)

	if txNonce != expectedNonce {
		logger.WithField("nonce.expected", expectedNonce).
			WithField("nonce.got", txNonce).
			Warnf("invalid nonce")

		if nc.recovery.Recovering(job.UUID) > nc.maxRecovery {
			return errors.InternalError("reached max nonce recovery max")
		}

		// Envelope has not already been recovered
		if txNonce > expectedNonce && nc.recovery.Recovering(nonceKey) == 0 {
			job.InternalData.ExpectedNonce = strconv.FormatUint(expectedNonce, 10)
		} else if txNonce < expectedNonce {
			// If nonce is to low we remove any recovery signal in metadata (possibly coming from a prior execution)
			job.InternalData.ExpectedNonce = ""
		}

		nc.recovery.Recover(job.UUID)
		return errors.InvalidNonceWarning("tx nonce expected %d, got %d", expectedNonce, txNonce)
	}

	return nil
}

func (nc *checker) OnFailure(ctx context.Context, job *entities.Job, jobErr error) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	if _, skip := nc.shouldSkip(job); skip {
		return nil
	}

	logger.Debug("checking job nonce on failure")

	// TODO: update EthClient to process and standardize nonce too low errors
	if !strings.Contains(jobErr.Error(), "nonce too low") &&
		!strings.Contains(jobErr.Error(), "Nonce too low") &&
		!strings.Contains(jobErr.Error(), "Incorrect nonce") {
		return nil
	}

	nonceKey := partitionKey(job)
	logger.Warnf("chain responded with invalid nonce error")
	if nc.recovery.Recovering(job.UUID) > nc.maxRecovery {
		return errors.InternalError("reached max nonce recovery max")
	}

	expectedNonce, err := nc.fetchNonceFromChain(ctx, job)
	if err != nil {
		logger.WithError(err).Error(fetchNonceErr)
		return err
	}

	logger.WithField("nonce.pending", expectedNonce).Debug("recalibrating nonce")
	err = nc.nonce.SetLastSent(nonceKey, expectedNonce-1)
	if err != nil {
		errMsg := "cannot set lastSent nonce"
		logger.WithError(err).Error(errMsg)
		return err
	}

	nc.recovery.Recover(job.UUID)

	return errors.InvalidNonceWarning(jobErr.Error())
}

func (nc *checker) OnSuccess(ctx context.Context, job *entities.Job) error {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	if _, skip := nc.shouldSkip(job); skip {
		return nil
	}

	logger.Debug("checking job nonce on success")

	nonceKey := partitionKey(job)

	nc.recovery.Recovered(job.UUID)
	txNonce, _ := strconv.ParseUint(job.Transaction.Nonce, 10, 32)
	err := nc.nonce.SetLastSent(nonceKey, txNonce)
	if err != nil {
		errMsg := "could not store last sent nonce"
		logger.WithError(err).Error(errMsg)
		return err
	}

	return nil
}

func (nc *checker) fetchNonceFromChain(ctx context.Context, job *entities.Job) (n uint64, err error) {
	url := fmt.Sprintf("%s/%s", nc.chainRegistryURL, job.ChainUUID)
	fromAddr := ethcommon.HexToAddress(job.Transaction.From)

	switch {
	case job.Type == tx.JobType_ETH_ORION_EEA_TX.String() && job.Transaction.PrivacyGroupID != "":
		n, err = nc.ethClient.PrivNonce(ctx, url, fromAddr,
			job.Transaction.PrivacyGroupID)
	case job.Type == tx.JobType_ETH_ORION_EEA_TX.String() && len(job.Transaction.PrivateFor) > 0:
		n, err = nc.ethClient.PrivEEANonce(ctx, url, fromAddr,
			job.Transaction.PrivateFrom, job.Transaction.PrivateFor)
	default:
		n, err = nc.ethClient.PendingNonceAt(ctx, url, fromAddr)
	}

	return
}

func (nc *checker) shouldSkip(job *entities.Job) (string, bool) {
	if job.InternalData.OneTimeKey {
		return "job is using one time key", true
	}
	if job.InternalData.ParentJobUUID != "" {
		return "job is a child", true
	}

	return "", false
}

func partitionKey(job *entities.Job) string {
	// Return empty partition key for raw tx and one time key tx
	// Not able to format a correct partition key if From or ChainID are not set. In that case return empty partition key
	if job.Transaction.From == "" || job.InternalData.ChainID == "" {
		return ""
	}

	fromAddr := job.Transaction.From
	chainID := job.InternalData.ChainID
	switch {
	case job.Type == tx.JobType_ETH_ORION_EEA_TX.String() && job.Transaction.PrivacyGroupID != "":
		return fmt.Sprintf("%v@orion-%v@%v", fromAddr, job.Transaction.PrivacyGroupID, chainID)
	case job.Type == tx.JobType_ETH_ORION_EEA_TX.String() && len(job.Transaction.PrivateFor) > 0:
		l := append(job.Transaction.PrivateFor, job.Transaction.PrivateFrom)
		sort.Strings(l)
		h := md5.New()
		_, _ = h.Write([]byte(strings.Join(l, "-")))
		return fmt.Sprintf("%v@orion-%v@%v", fromAddr, fmt.Sprintf("%x", h.Sum(nil)), chainID)
	default:
		return fmt.Sprintf("%v@%v", fromAddr, chainID)
	}
}
