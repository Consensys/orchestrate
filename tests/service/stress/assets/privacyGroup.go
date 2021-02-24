package assets

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/ethclient"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/utils"
)

var privacyGroupsCtxKey ctxKey = "besuPrivacyGroups"

type PrivacyGroup struct {
	ID    string
	Nodes []string
}

func CreatePrivateGroup(ctx context.Context, ec ethclient.EEAClient, chainURL string, myselfNodeAddress, otherNodeAddress []string) (context.Context, error) {
	logger := log.FromContext(ctx).WithField("url", chainURL)
	logger.Debug("creating new privacy group")

	size := utils.RandIntRange(1, len(otherNodeAddress))
	nMyselfNode := utils.RandInt(len(myselfNodeAddress))
	nodes := append(otherNodeAddress[0:size], myselfNodeAddress[nMyselfNode])

	logger = logger.WithField("address", nodes)
	privateGroupID, err := ec.PrivCreatePrivacyGroup(ctx, chainURL, nodes)
	if err != nil {
		errMsg := "failed to register besu privacy group"
		logger.WithError(err).Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	logger.WithField("id", privateGroupID).Info("privacy group has been registered")
	return contextWithPrivacyGroups(ctx, append(ContextPrivacyGroups(ctx),
		PrivacyGroup{ID: privateGroupID, Nodes: nodes})), nil
}

func contextWithPrivacyGroups(ctx context.Context, groups []PrivacyGroup) context.Context {
	return context.WithValue(ctx, privacyGroupsCtxKey, groups)
}

func ContextPrivacyGroups(ctx context.Context) []PrivacyGroup {
	if v, ok := ctx.Value(privacyGroupsCtxKey).([]PrivacyGroup); ok {
		return v
	}
	return []PrivacyGroup{}
}
