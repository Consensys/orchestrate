package contractregistry

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

func TestApp(t *testing.T) {
	ctrlr := gomock.NewController(t)
	defer ctrlr.Finish()

	jwtChecker := mockauth.NewMockChecker(ctrlr)
	keyChecker := mockauth.NewMockChecker(ctrlr)

	cfg := NewConfig(viper.New())

	_, err := New(
		cfg,
		postgres.GetManager(),
		jwtChecker, keyChecker,
	)
	assert.NoError(t, err, "Creating App should not error")
}
