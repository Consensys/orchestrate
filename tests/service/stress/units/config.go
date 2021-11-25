package units

import (
	"context"
	"time"

	"github.com/consensys/orchestrate/tests/service/stress/assets"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type WorkloadConfig struct {
	accounts               []ethcommon.Address
	chains                 []assets.Chain
	artifacts              []string
	privacyGroups          []assets.PrivacyGroup
	waitForEnvelopeTimeout time.Duration
}

func NewWorkloadConfig(ctx context.Context, waitForEnvelopeTimeout time.Duration) *WorkloadConfig {
	return &WorkloadConfig{
		accounts:               assets.ContextAccounts(ctx),
		chains:                 assets.ContextChains(ctx),
		artifacts:              assets.ContextArtifacts(ctx),
		privacyGroups:          assets.ContextPrivacyGroups(ctx),
		waitForEnvelopeTimeout: waitForEnvelopeTimeout,
	}
}
