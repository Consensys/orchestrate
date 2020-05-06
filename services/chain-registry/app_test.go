package chainregistry

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	mockethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/mock"
)

func TestApp(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	jwtChecker := mockauth.NewMockChecker(ctrlr)
	keyChecker := mockauth.NewMockChecker(ctrlr)
	ec := mockethclient.NewMockClient(ctrlr)

	cfg := NewConfig(viper.New())
	cfg.Store.Type = "postgres"

	_, err := New(
		cfg,
		postgres.GetManager(),
		ec,
		jwtChecker, keyChecker,
	)
	assert.NoError(t, err, "Creating App should not error")
}
