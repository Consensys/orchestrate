package testutils

import (
	"fmt"

	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/identitymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func FakeCreateIdentityRequest() *types.CreateIdentityRequest {
	return &types.CreateIdentityRequest{
		Alias: fmt.Sprintf("Alias_%s", utils.RandomString(5)),
		Attributes: map[string]string{
			"attr1": "val1",
			"attr2": "val2",
		},
	}
}

func FakeImportIdentityRequest() *types.ImportIdentityRequest {
	return &types.ImportIdentityRequest{
		Alias:      fmt.Sprintf("Alias_%s", utils.RandomString(5)),
		PrivateKey: "66232652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2D",
		Attributes: map[string]string{
			"attr1": "val1",
			"attr2": "val2",
		},
	}
}
