// +build unit

package envelopestore

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/mock"
)

func TestApp(t *testing.T) {
	ctlr := gomock.NewController(t)
	defer ctlr.Finish()

	cfg := app.DefaultConfig()
	appli, err := New(
		cfg,
		mockauth.NewMockChecker(ctlr),
		mockauth.NewMockChecker(ctlr),
		false,
		mockstore.NewMockEnvelopeStoreServer(ctlr),
		logrus.New(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, appli, "App should have been created")
}
