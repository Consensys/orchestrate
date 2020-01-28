package multitenancy

import (
	"context"
	"strings"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

var sep = "@"

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

	tenantID := TenantIDFromContext(ctx)
	if tenantID == "" {
		return "", errors.NotFoundError("not able to retrieve the tenant UUID: The tenant_id is not present in the Context")
	}

	newKey := tenantID + key

	return newKey, nil
}

func SplitTenant(key string) (context.Context, string, error) {
	slicePkey := strings.Split(key, sep)
	switch len(slicePkey) {
	case 1:
		return context.Background(), key, nil
	case 2:
		ctx := WithTenantID(context.Background(), slicePkey[1])
		return ctx, slicePkey[0], nil
	default:
		return nil, "", errors.InvalidFormatError("The key have more than one separator as " + sep)
	}
}
