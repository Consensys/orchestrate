package testutils

import (
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/identitymanager"
)

func FakeCreateIdentityRequest() *types.CreateIdentityRequest {
	return &types.CreateIdentityRequest{
		Alias: "MyPersonalAccount",
		Attributes: map[string]string{
			"attr1": "val1",
			"attr2": "val2",
		},
	}
}
