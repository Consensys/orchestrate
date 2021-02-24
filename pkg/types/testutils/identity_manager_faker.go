package testutils

import (
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/types/api"

	"github.com/ConsenSys/orchestrate/pkg/utils"
)

func FakeCreateAccountRequest() *api.CreateAccountRequest {
	return &api.CreateAccountRequest{
		Alias: fmt.Sprintf("Alias_%s", utils.RandString(5)),
		Attributes: map[string]string{
			"attr1": "val1",
			"attr2": "val2",
		},
	}
}

func FakeImportAccountRequest() *api.ImportAccountRequest {
	return &api.ImportAccountRequest{
		Alias:      fmt.Sprintf("Alias_%s", utils.RandString(5)),
		PrivateKey: "a93e498896143c02fdf42b9b69bdcf4aebcedc8d45851c33f8ae86057e7c4a90",
		Attributes: map[string]string{
			"attr1": "val1",
			"attr2": "val2",
		},
	}
}

func FakeUpdateAccountRequest() *api.UpdateAccountRequest {
	return &api.UpdateAccountRequest{
		Alias: fmt.Sprintf("Alias_%s", utils.RandString(5)),
		Attributes: map[string]string{
			"attr3": "val3",
			"attr4": "val4",
		},
	}
}
