package testutils

import (
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func FakeCreateAccountRequest() *api.CreateAccountRequest {
	return &api.CreateAccountRequest{
		Alias: fmt.Sprintf("Alias_%s", utils.RandomString(5)),
		Attributes: map[string]string{
			"attr1": "val1",
			"attr2": "val2",
		},
	}
}

func FakeImportAccountRequest() *api.ImportAccountRequest {
	return &api.ImportAccountRequest{
		Alias:      fmt.Sprintf("Alias_%s", utils.RandomString(5)),
		PrivateKey: "a93e498896143c02fdf42b9b69bdcf4aebcedc8d45851c33f8ae86057e7c4a90",
		Attributes: map[string]string{
			"attr1": "val1",
			"attr2": "val2",
		},
	}
}

func FakeUpdateAccountRequest() *api.UpdateAccountRequest {
	return &api.UpdateAccountRequest{
		Alias: fmt.Sprintf("Alias_%s", utils.RandomString(5)),
		Attributes: map[string]string{
			"attr3": "val3",
			"attr4": "val4",
		},
	}
}
