// +build integration

package integrationtests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
	"math/rand"
	"testing"
)

// faucetTestSuite is a test suite for faucet API
type faucetTestSuite struct {
	suite.Suite
	baseURL string
	client  client.FaucetClient
	env     *IntegrationEnvironment
}

func (s *faucetTestSuite) SetupSuite() {
	s.client = client.DialWithDefaultOptions(s.baseURL)
}

func (s *faucetTestSuite) TestChainRegistry_FaucetHappyFlow() {
	ctx := context.Background()
	faucetName := fmt.Sprintf("TestFaucet%d", rand.Intn(1000))
	faucetNameTwo := fmt.Sprintf("TestFaucet%d", rand.Intn(1000))
	var faucetUUID string

	s.T().Run("should register a new faucet", func(t *testing.T) {
		faucet := models.Faucet{
			Name:            faucetName,
			ChainRule:       "public/",
			CreditorAccount: "0x8dd688660ec0BaBD0B8a2f2DE3232645F73cC5eb",
			MaxBalance:      "1000",
			Amount:          "1000",
			Cooldown:        "1s",
		}
		resp, err := s.client.RegisterFaucet(ctx, &faucet)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.NotEmpty(t, resp.UUID)
		faucetUUID = resp.UUID
	})

	s.T().Run("should update registered faucet by UUID", func(t *testing.T) {
		faucet := &models.Faucet{
			Name: faucetNameTwo,
		}

		_, err := s.client.UpdateFaucetByUUID(ctx, faucetUUID, faucet)
		assert.NoError(t, err)

		faucet, err = s.client.GetFaucetByUUID(ctx, faucetUUID)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, faucet.Name, faucetNameTwo)
	})

	s.T().Run("should delete registered faucet by UUID", func(t *testing.T) {
		err := s.client.DeleteFaucetByUUID(ctx, faucetUUID)
		assert.NoError(t, err)

		_, err = s.client.GetFaucetByUUID(ctx, faucetUUID)
		assert.True(t, errors.IsNotFoundError(err))
	})
}
