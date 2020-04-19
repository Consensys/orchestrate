// +build unit

package envelopestore

import (
	"testing"

	"github.com/golang/mock/gomock"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	mockstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/controllers/mock"
)

func TestApp(t *testing.T) {
	ctlr := gomock.NewController(t)
	defer ctlr.Finish()

	cfg := NewConfigFromViper(viper.GetViper())
	appli, err := newApplication(
		&cfg,
		mockauth.NewMockChecker(ctlr),
		mockauth.NewMockChecker(ctlr),
		mockstore.NewMockEnvelopeStoreServer(ctlr),
		logrus.New(),
		prom.NewRegistry(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, appli, "App should have been created")
}
