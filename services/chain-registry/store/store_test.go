// +build unit

package store_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

func TestImportChains(t *testing.T) {
	testSuite := []struct {
		name   string
		chains []string
	}{
		{
			"import chains",
			[]string{`{"name":"noError"}`},
		},
		{
			"import chains with error",
			[]string{`{"name":"error"}`},
		},
		{
			"import chains with unknown field",
			[]string{`{"unknown":"error"}`},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockStore := mockstore.NewMockChainRegistryStore(mockCtrl)
	mockStore.EXPECT().RegisterChain(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, chain *types.Chain) error {
			if chain.Name == "error" {
				return fmt.Errorf("error")
			}
			return nil
		}).AnyTimes()
	mockStore.EXPECT().UpdateChainByName(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	for _, test := range testSuite {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			store.ImportChains(context.Background(), mockStore, test.chains)
		})
	}
}
