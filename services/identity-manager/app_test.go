// +build unit

package identitymanager

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
)

func TestApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtChecker := mockauth.NewMockChecker(ctrl)
	keyChecker := mockauth.NewMockChecker(ctrl)

	cfg := NewConfig(viper.New())
	cfg.Store.Type = "postgres"

	_, err := NewIdentityManager(cfg, postgres.GetManager(), jwtChecker, keyChecker)
	assert.NoError(t, err, "Creating App should not error")
}
