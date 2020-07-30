package multitenancy

import (
	"context"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type KeyBuilder struct {
	multitenancy bool
}

func New(multiTenancyEnabled bool) *KeyBuilder {
	return &KeyBuilder{
		multitenancy: multiTenancyEnabled,
	}
}

func (k *KeyBuilder) BuildKey(ctx context.Context, key string) (string, error) {
	tenant := TenantIDFromContext(ctx)

	return k.BuildKeyWithTenant(tenant, key), nil
}

func (k *KeyBuilder) BuildKeyWithTenant(tenantID, key string) string {
	return tenantID + key
}

func SplitTenant(key string) (context.Context, string, error) {
	slicePkey := strings.Split(key, "@")
	switch len(slicePkey) {
	case 1:
		return context.Background(), key, nil
	case 2:
		ctx := WithTenantID(context.Background(), slicePkey[1])
		return ctx, slicePkey[0], nil
	default:
		return nil, "", errors.InvalidFormatError("invalid key %v", key)
	}
}
