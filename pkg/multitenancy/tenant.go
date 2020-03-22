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
	if !k.multitenancy {
		return key, nil
	}

	return TenantIDFromContext(ctx) + key, nil
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
