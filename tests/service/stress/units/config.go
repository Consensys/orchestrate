package units

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/stress/assets"
)

type WorkloadConfig struct {
	accounts               []string
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
