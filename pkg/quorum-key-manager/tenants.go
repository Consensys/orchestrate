package quorumkeymanager

import (
	"context"
	"strings"

	qkmclient "github.com/ConsenSys/orchestrate/pkg/quorum-key-manager/client"
	"github.com/ConsenSys/orchestrate/pkg/utils"
)

const (
	TagIDAllowedTenants        = "tenants"
	TagSeparatorAllowedTenants = ","
)

func IsTenantAllowed(ctx context.Context, client qkmclient.Eth1Client, tenants []string, storeName, address string) (bool, error) {
	acc, err := client.GetEth1Account(ctx, storeName, address)
	if err != nil {
		return false, err
	}

	allowedTenants := strings.Split(acc.Tags[TagIDAllowedTenants], TagSeparatorAllowedTenants)
	return len(utils.ArrayIntersection(tenants, allowedTenants).([]interface{})) > 0, nil
}
