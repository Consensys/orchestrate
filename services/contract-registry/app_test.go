// +build unit

package contractregistry

import (
	"testing"

	"github.com/golang/mock/gomock"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	mockauth "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/mock"
	mockregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service/mock"
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
		mockregistry.NewMockContractRegistryServer(ctlr),
		logrus.New(),
		prom.NewRegistry(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, appli, "App should have been created")
}
